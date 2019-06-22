package opr

import (
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/FactomProject/btcutil/base58"
	"github.com/FactomProject/factom"
	"github.com/pegnet/PegNet/support"
	config "github.com/zpatrick/go-config"
)

// Compute the average answer for the price of each token reported
func Avg(list []*OraclePriceRecord) (avg [20]float64) {
	// Sum up all the prices
	for _, opr := range list {
		tokens := opr.GetTokens()
		for i, price := range tokens {
			avg[i] += price.value
		}
	}
	// Then divide the prices by the number of OraclePriceRecord records.  Two steps is actually faster
	// than doing everything in one loop (one divide for every asset rather than one divide
	// for every asset * number of OraclePriceRecords)
	numList := float64(len(list))
	for i := range avg {
		avg[i] = avg[i] / numList / 100000000
	}
	return
}

// Given the average answers across a set of tokens, grade the opr
func CalculateGrade(avg [20]float64, opr *OraclePriceRecord) float64 {
	tokens := opr.GetTokens()
	for i, v := range tokens {
		d := v.value - avg[i]           // compute the difference from the average
		opr.Grade = opr.Grade + d*d*d*d // the grade is the sum of the squares of the differences
	}
	return opr.Grade
}

// Given a list of OraclePriceRecord, figure out which 10 should be paid, and in what order
func GradeBlock(list []*OraclePriceRecord) (tobepaid []*OraclePriceRecord, sortedlist []*OraclePriceRecord) {

	if len(list) < 10 {
		return nil, nil
	}

	last := len(list)
	// Throw away all the entries but the top 50 in difficulty
	// bubble sort because I am lazy.  Could be replaced with about anything
	for j := 0; j < len(list)-1; j++ {
		for k := 0; k < len(list)-j-1; k++ {
			d1 := list[k].Difficulty
			d2 := list[k+1].Difficulty
			if d1 == 0 || d2 == 0 {
				//panic("Should not be here")
			}
			if d1 < d2 { // sort the smallest difficulty to the end of the list
				list[k], list[k+1] = list[k+1], list[k]
			}
		}
	}
	if len(list) > 50 {
		last = 50
	}
	// Go through and throw away entries that are outside the average or on a tie, have the worst difficulty
	// until we are only left with 10 entries to reward
	for i := last; i >= 10; i-- {
		avg := Avg(list[:i])
		for j := 0; j < i; j++ {
			CalculateGrade(avg, list[j])
		}
		// bubble sort the worst grade to the end of the list. Note that this is nearly sorted data, so
		// a bubble sort with a short circuit is pretty darn good sort.
		for j := 0; j < i-1; j++ {
			cont := false                // If we can get through a pass with no swaps, we are done.
			for k := 0; k < i-j-1; k++ { // yes, yes I know we can get 2 or 3 x better speed playing with indexes
				if list[k].Grade > list[k+1].Grade { // bit it is tricky.  This is good enough.
					list[k], list[k+1] = list[k+1], list[k] // sort first by the grade.
					cont = true                             // any swap means we continue to loop
				} else if list[k].Grade == list[k+1].Grade { // break ties with PoW.  Where data is being shared
					if list[k].Difficulty < list[k+1].Difficulty { // we will have ties.
						//list[k], list[k+1] = list[k+1], list[k]
						cont = true // any swap means we continue to loop
					}
				}
			}
			if !cont { // If we made a pass without any swaps, we are done.
				break
			}
		}
	}
	tobepaid = append(tobepaid, list[:10]...)
	return tobepaid, list
}

type OPRBlock struct {
	OPRs []*OraclePriceRecord
	Dbht int64
}

var OPRBlocks []*OPRBlock
var EBMutex sync.Mutex

// Get the OPR Records at a given dbht
func GetEntryBlocks(config *config.Config) {
	EBMutex.Lock()
	defer EBMutex.Unlock()

	p, err := config.String("Miner.Protocol")
	check(err)
	n, err := config.String("Miner.Network")
	check(err)
	opr := [][]byte{[]byte(p), []byte(n), []byte("Oracle Price Records")}
	heb, err := factom.GetChainHead(hex.EncodeToString(support.ComputeChainIDFromFields(opr)))
	check(err)
	eb, err := factom.GetEBlock(heb)
	check(err)

	// A temp list of candidate oprblocks to evaluate to see if they fit nicely together
	// Because we go from the head of the chain backwards to collect them, they have to be
	// collected before I can then validate them forward from the highest valid OPR block
	// I have found.
	var oprblocks []*OPRBlock
	// For each entryblock in the Oracle Price Records chain
	// Get all the valid OPRs and put them in  a new OPRBlock structure
	for eb != nil && (len(OPRBlocks) == 0 ||
		eb.Header.DBHeight > OPRBlocks[len(OPRBlocks)-1].Dbht) {

		// Go through the Entry Block and collect all the valid OPR records
		if len(eb.EntryList) > 10 {
			oprblk := new(OPRBlock)
			oprblk.Dbht = eb.Header.DBHeight
			for _, ebentry := range eb.EntryList {
				entry, err := factom.GetEntry(ebentry.EntryHash)
				check(err)

				// Do some quick collecting of data and checks of the entry.
				// Can only have one ExtID which must be the nonce for the entry
				if len(entry.ExtIDs) != 1 || len(entry.ExtIDs[0]) != 32 {
					continue // keep looking if the entry has more than one extid or it isn't 32 bytes
				}

				// Okay, it looks sort of okay.  Lets unmarshal the JSON
				opr := new(OraclePriceRecord)
				if err := json.Unmarshal(entry.Content, opr); err != nil {
					continue // Doesn't unmarshal, then it isn't valid for sure.  Continue on.
				}

				// Run some basic checks on the values.  If they don't check out, then ignore the entry
				if !opr.Validate(config) {
					continue
				}

				// Looking good.  Go ahead and compute the OPRHash
				opr.OPRHash = LX.Hash(entry.Content) // Save the OPRHash
				opr.Entry = entry                    // Compute the Entry Hash
				eh, err := hex.DecodeString(ebentry.EntryHash)
				if err != nil {
					continue
				}
				opr.EntryHash = base58.Encode(eh) // Encode to base58

				// Okay, mostly good.  Add to our candidate list
				oprblk.OPRs = append(oprblk.OPRs, opr)

			}
			// If we have 10 canidates, then lets add them up.
			if len(oprblk.OPRs) >= 10 {
				oprblocks = append(oprblocks, oprblk)
			}
		}
		// At this point, the oprblk has all the valid OPRs. Make sure we have enough.
		// sorted list of winners.

		neb, err := factom.GetEBlock(eb.Header.PrevKeyMR)
		if err != nil {
			break
		}
		eb = neb
	}

	// Take the reverse ordered oprblocks, from last to first.  Validate all the winners are
	// the right winners.  Replace the generally correct OPR list in the oprblock with the
	// list of winners.  These should be the winners of the next block, which lucky enough is
	// the next block we are going to process.
	// Ignore oprblocks that don't get 10 winners.
	for i := len(oprblocks) - 1; i >= 0; i-- { // Okay, go through these backwards
		prevOPRBlock := GetPreviousOPRs(int32(oprblocks[i].Dbht)) // Get the previous OPRBlock
		var validOPRs []*OraclePriceRecord                        // Collect the valid OPRPriceRecords here
		for _, opr := range oprblocks[i].OPRs {                   // Go through this block
			for j, eh := range opr.WinningPreviousOPR { // Make sure the winning records are valid
				if (prevOPRBlock == nil && eh != "") ||
					(prevOPRBlock != nil && eh != prevOPRBlock[0].WinningPreviousOPR[j]) {
					continue
				}
				opr.Difficulty = opr.ComputeDifficulty(opr.Entry.ExtIDs[0])
			}
			validOPRs = append(validOPRs, opr) // Add to my valid list if all the winners are right
		}
		if len(validOPRs) < 10 { // Make sure we have at least 10 valid OPRs,
			continue // and leave if we don't.
		}
		winners, _ := GradeBlock(validOPRs)
		oprblocks[i].OPRs = winners
		OPRBlocks = append(OPRBlocks, oprblocks[i])
	}

	return
}

// GetPreviousOPRs()
// So what they are asking for here is the previous winning blocks. In our list, we have graded and ordered
// the OPRs, so just go through the list and return the highest dbht less than the one asked for.
// Returns nil if the dbht is the first dbht in the chain.
func GetPreviousOPRs(dbht int32) []*OraclePriceRecord {
	for i := len(OPRBlocks) - 1; i >= 0; i-- {
		if OPRBlocks[i].Dbht < int64(dbht) {
			return OPRBlocks[i].OPRs
		}
	}
	return nil
}
