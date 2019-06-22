package support

import (
	"fmt"
	"strings"

	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/FactomProject/btcutil/base58"
	"bytes"
)

type NetworkType int

const (
	INVALID NetworkType = iota + 1
	MAIN_NETWORK
	TEST_NETWORK
)

var AssetNames = []string{
	"PNT",
	"USD",
	"EUR",
	"JPY",
	"GBP",
	"CAD",
	"CHF",
	"INR",
	"SGD",
	"CNY",
	"HKD",
	"XAU",
	"XAG",
	"XPD",
	"XPT",
	"XBT",
	"ETH",
	"LTC",
	"XBC",
	"FCT",
}

var PegAssetNames []string

var TestPegAssetNames []string

func init() {
	for _,asset := range AssetNames {
		PegAssetNames = append(PegAssetNames, "p"+asset)
		TestPegAssetNames = append(TestPegAssetNames, "t"+asset)
	}
}

func PullValue(line string, howMany int) string {
	i := 0
	//fmt.Println(line)
	var pos int
	for i < howMany {
		//find the end of the howmany-th tag
		pos = strings.Index(line, ">")
		line = line[pos+1:]
		//fmt.Println(line)
		i = i + 1
	}
	//fmt.Println("line:", line)
	pos = strings.Index(line, "<")
	//fmt.Println("POS:", pos)
	line = line[0:pos]
	//fmt.Println(line)
	return line
}

// CheckPrefix()
// Check the prefix for either network type.
func CheckPrefix(network NetworkType, name string) bool {
	if network == MAIN_NETWORK {
		for _, v := range PegAssetNames {
			if v == name {
				return true
			}
		}
	} else {
		for _, v := range TestPegAssetNames {
			if v == name {
				return true
			}
		}
	}
	return false

}

// ConvertRawAddrToPeg()
// Converts a raw RCD1 address into a wallet friendly address that can be used to
// convert assets, check balances, and send tokens.  While the underlying private key can be
// used to hold Factoids or any token in the PegNet, users need addresses that create a
// barrier to mistakes that can lead to sending the wrong tokens to the wrong addresses
func ConvertRawAddrToPeg(network NetworkType, prefix string, adr []byte) (string, error) {

	// Make sure the prefix is valid.
	if !CheckPrefix(network, prefix) {
		return "", errors.New(prefix + " is not a valid PegNet prefix")
	}

	h := sha256.Sum256([]byte(append([]byte(prefix),adr...)))
	hash := sha256.Sum256(h[:])

	// Append the prefix to the base 58 representation of the raw address
	b58 := prefix +"_" + base58.Encode( append(adr, hash[:4]...))

	return b58, nil
}

// ConvertPegTAddrToRaw()
// Convert a human/wallet address to the raw underlying address.  Verifies the checksum and
// the validity of the prefix.  Returns the prefix, the raw address, and error.
//
func ConvertPegTAddrToRaw(network NetworkType, adr string) (prefix string, rawAdr []byte, err error) {
	adrLen := len(adr)
	if adrLen < 42 || len(adr) > 56 {
		return "", nil, errors.New(
			fmt.Sprintf("valid pegNet token addresses are 44 to 56 characters in length. len(adr)=%d ", adrLen))
	}

	prefix = adr[:4]
	if !CheckPrefix(network, prefix) {
		return "", nil, errors.New(prefix + " is not a valid PegNet prefix")
	}

	b58 := adr[5:]
	raw := base58.Decode(b58)
	if len(raw) == 0 {
		return "", nil, errors.New("invalid base58 encoding")
	}
	rawAdr = raw[:len(raw)-4]
	chksum := raw[len(raw)-4:]

	hash := sha256.Sum256(append([]byte(prefix), rawAdr...))
	hash = sha256.Sum256(hash[:])
	if bytes.Compare(hash[:4], chksum)!=0 {
		return "", nil, errors.New("checksum failure")
	}

	return prefix, rawAdr, nil

}

// PegTAdrIsValid()
// Check that the given human/wallet PegNet address is valid.
func PegTAdrIsValid(network NetworkType, adr string) error {
	_, _, err := ConvertPegTAddrToRaw(network, adr)
	return err
}

// RandomByteSliceOfLen()
// Returns a random set of bytes of a given length
func RandomByteSliceOfLen(sliceLen int) []byte {
	if sliceLen <= 0 {
		return nil
	}
	answer := make([]byte, sliceLen)
	_, err := rand.Read(answer)
	if err != nil {
		return nil
	}
	return answer
}

//  Convert Factoid and Entry Credit addresses to their more user
//  friendly and human readable formats.
//
//  Creates the binary form.  Just needs the conversion to base58
//  for display.
func ConvertFctAddressToUser(addr []byte) string {
	dat := make([]byte, 0, 64)
	dat = append(dat, 0x5f, 0xb1)
	dat = append(dat, addr...)
	hash := sha256.Sum256(dat)
	sha256d := sha256.Sum256(hash[:])
	userd := []byte{0x5f, 0xb1}
	userd = append(userd, addr...)
	userd = append(userd, sha256d[:4]...)
	return base58.Encode(userd)
}

//  Convert Factoid and Entry Credit addresses to their more user
//  friendly and human readable formats.
//
//  Creates the binary form.  Just needs the conversion to base58
//  for display.
func ConvertECAddressToUser(addr []byte) string {
	dat := make([]byte, 0, 64)
	dat = append(dat, 0x59, 0x2a)
	dat = append(dat, addr...)
	hash := sha256.Sum256(dat)
	sha256d := sha256.Sum256(hash[:])
	userd := []byte{0x59, 0x2a}
	userd = append(userd, addr...)
	userd = append(userd, sha256d[:4]...)
	return base58.Encode(userd)
}


// Convert a User facing Factoid or Entry Credit address
// or their Private Key representations
// to the regular form.  Note validation must be done
// separately!
func ConvertFctECUserStrToAddress(userFAddr string) string {
	v := base58.Decode(userFAddr)
	return hex.EncodeToString(v[2:34])
}


// Convert a User facing FCT address to all of its PegNet
// asset token User facing forms.
func ConvertUserFctToUserPegNetAssets(userFctAddr string)(assets[]string, err error){
	raw := base58.Decode(userFctAddr)[2:34]
	cvt :=func(network NetworkType, asset string) (passet string) {
		passet, err = ConvertRawAddrToPeg(network, asset, raw)
		if err != nil {
			panic(err)
		}
		return passet
	}

	for _, asset := range AssetNames {
		assets = append(assets, cvt(MAIN_NETWORK, "p"+asset))
		assets = append(assets, cvt(TEST_NETWORK, "t"+asset))
	}

	return assets, nil
}