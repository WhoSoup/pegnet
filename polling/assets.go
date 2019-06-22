// Copyright (c) of parts are held by the various contributors (see the CLA)
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
package polling

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/zpatrick/go-config"
	"strconv"
	"sync"
	"time"
)

const qlimit = 580 // Limit queries to once just shy of 10 minutes (600 seconds)

type PegAssets struct {
	PNT PegItems
	USD PegItems
	EUR PegItems
	JPY PegItems
	GBP PegItems
	CAD PegItems
	CHF PegItems
	INR PegItems
	SGD PegItems
	CNY PegItems
	HKD PegItems
	XAU PegItems
	XAG PegItems
	XPD PegItems
	XPT PegItems
	XBT PegItems
	ETH PegItems
	LTC PegItems
	XBC PegItems
	FCT PegItems
}

func (p *PegAssets) Clone() PegAssets {
	np := new(PegAssets)
	np.PNT = p.PNT.Clone()
	np.USD = p.USD.Clone()
	np.EUR = p.EUR.Clone()
	np.JPY = p.JPY.Clone()
	np.GBP = p.GBP.Clone()
	np.CAD = p.CAD.Clone()
	np.CHF = p.CHF.Clone()
	np.INR = p.INR.Clone()
	np.SGD = p.SGD.Clone()
	np.CNY = p.CNY.Clone()
	np.HKD = p.HKD.Clone()
	np.XAU = p.XAU.Clone()
	np.XAG = p.XAG.Clone()
	np.XPD = p.XPD.Clone()
	np.XPT = p.XPT.Clone()
	np.XBT = p.XBT.Clone()
	np.ETH = p.ETH.Clone()
	np.LTC = p.LTC.Clone()
	np.XBC = p.XBC.Clone()
	np.FCT = p.FCT.Clone()
	return *np
}

type PegItems struct {
	Value float64
	When  string
}

func (p *PegItems) Clone() PegItems {
	np := new(PegItems)
	np.Value = p.Value
	np.When = p.When
	return *np
}

var lastMutex sync.Mutex
var lastAnswer PegAssets //
var lastTime int64       // In seconds

func Round(v float64) float64 {
	return float64(int64(v*10000)) / 10000
}

func PullPEGAssets(config *config.Config) (pa PegAssets) {

	// Prevent pounding of external APIs
	lastMutex.Lock()
	defer lastMutex.Unlock()
	now := time.Now().Unix()
	delta := now - lastTime
	if delta < qlimit && lastTime != 0 {
		pa := lastAnswer.Clone()
		return pa
	}

	lastTime = now
	fmt.Println("Make a call to get data. Seconds since last call:", delta)
	var Peg PegAssets
	// digital currencies
	CoinCapResponseBytes, err := CallCoinCap(config)
	//fmt.Println(string(CoinCapResponseBytes))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		var CoinCapValues CoinCapResponse
		err = json.Unmarshal(CoinCapResponseBytes, &CoinCapValues)
		for _, currency := range CoinCapValues.Data {
			//fmt.Println(currency.Symbol + "-" + currency.PriceUSD)
			if currency.Symbol == "XBT" || currency.Symbol == "BTC" {
				Peg.XBT.Value, err = strconv.ParseFloat(currency.PriceUSD, 64)
				Peg.XBT.Value = Round(Peg.XBT.Value)
				if err != nil {
					continue
				}
				Peg.XBT.When = string(CoinCapValues.Timestamp)
			} else if currency.Symbol == "ETH" {
				Peg.ETH.Value, err = strconv.ParseFloat(currency.PriceUSD, 64)
				Peg.ETH.Value = Round(Peg.ETH.Value)
				if err != nil {
					continue
				}
				Peg.ETH.When = string(CoinCapValues.Timestamp)
			} else if currency.Symbol == "LTC" {
				Peg.LTC.Value, err = strconv.ParseFloat(currency.PriceUSD, 64)
				Peg.LTC.Value = Round(Peg.LTC.Value)
				if err != nil {
					continue
				}
				Peg.LTC.When = string(CoinCapValues.Timestamp)
			} else if currency.Symbol == "XBC" || currency.Symbol == "BCH" {
				Peg.XBC.Value, err = strconv.ParseFloat(currency.PriceUSD, 64)
				Peg.XBC.Value = Round(Peg.XBC.Value)
				if err != nil {
					continue
				}
				Peg.XBC.When = string(CoinCapValues.Timestamp)
			} else if currency.Symbol == "FCT" {
				Peg.FCT.Value, err = strconv.ParseFloat(currency.PriceUSD, 64)
				Peg.FCT.Value = Round(Peg.FCT.Value)
				if err != nil {
					continue
				}
				Peg.FCT.When = string(CoinCapValues.Timestamp)
			}
		}
	}

	//fiat option 1.  terms uf use seem tighter
	// has fiat and digital
	// https://currencylayer.com/product  <-- pricing
	// $10 a month will let you pull 10,000 but it is updated once an hour
	// $40 a month is 100,000 updated every 10 minutes.
	// 20% discount for annual payment

	//	fmt.Println("API LAYER:")

	APILayerBytes, err := CallAPILayer(config)
	//fmt.Println(string(APILayerBytes))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		var APILayerResponse APILayerResponse
		err = json.Unmarshal(APILayerBytes, &APILayerResponse)

		Peg.USD.Value = Round(APILayerResponse.Quotes.USDUSD)
		Peg.USD.When = string(APILayerResponse.Timestamp)
		Peg.EUR.Value = Round(APILayerResponse.Quotes.USDEUR)
		Peg.EUR.When = string(APILayerResponse.Timestamp)
		Peg.JPY.Value = Round(APILayerResponse.Quotes.USDJPY)
		Peg.JPY.When = string(APILayerResponse.Timestamp)
		Peg.GBP.Value = Round(APILayerResponse.Quotes.USDGBP)
		Peg.GBP.When = string(APILayerResponse.Timestamp)
		Peg.CAD.Value = Round(APILayerResponse.Quotes.USDCAD)
		Peg.CAD.When = string(APILayerResponse.Timestamp)
		Peg.CHF.Value = Round(APILayerResponse.Quotes.USDCHF)
		Peg.CHF.When = string(APILayerResponse.Timestamp)
		Peg.INR.Value = Round(APILayerResponse.Quotes.USDINR)
		Peg.INR.When = string(APILayerResponse.Timestamp)
		Peg.SGD.Value = Round(APILayerResponse.Quotes.USDSGD)
		Peg.SGD.When = string(APILayerResponse.Timestamp)
		Peg.CNY.Value = Round(APILayerResponse.Quotes.USDCNY)
		Peg.CNY.When = string(APILayerResponse.Timestamp)
		Peg.HKD.Value = Round(APILayerResponse.Quotes.USDHKD)
		Peg.HKD.When = string(APILayerResponse.Timestamp)

	}

	KitcoResponse, err := CallKitcoWeb()

	for i := 0; i < 10; i++ {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %d so retrying.  Error %v\n", i+1, err)
			time.Sleep(time.Second)
			KitcoResponse, err = CallKitcoWeb()
		} else {
			break //	os.Exit(1)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error, using old data.\n")
		pa := lastAnswer.Clone()
		return pa
	}
	//fmt.Println("KitcoResponse:", KitcoResponse)
	Peg.XAU.Value, err = strconv.ParseFloat(KitcoResponse.Silver.Bid, 64)
	Peg.XAU.When = KitcoResponse.Silver.Date
	Peg.XAG.Value, err = strconv.ParseFloat(KitcoResponse.Gold.Bid, 64)
	Peg.XAG.When = KitcoResponse.Gold.Date
	Peg.XPD.Value, err = strconv.ParseFloat(KitcoResponse.Palladium.Bid, 64)
	Peg.XPD.When = KitcoResponse.Palladium.Date
	Peg.XPT.Value, err = strconv.ParseFloat(KitcoResponse.Platinum.Bid, 64)
	Peg.XPT.When = KitcoResponse.Platinum.Date

	lastAnswer = Peg.Clone()
	return Peg
}

func (peg *PegAssets) FillPriceBytes() {
	byteVal := make([]byte, 160)
	nextStart := 0
	byteLength := 8
	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, uint64(peg.PNT.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.USD.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.EUR.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.JPY.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.GBP.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.CAD.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.CHF.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.INR.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.SGD.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.CNY.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.HKD.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.XAU.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.XAG.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.XPD.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.XPT.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.XBT.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.ETH.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.LTC.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.XBC.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
	nextStart = nextStart + byteLength
	binary.BigEndian.PutUint64(b, uint64(peg.FCT.Value))
	copy(byteVal[nextStart:nextStart+8], b[:])
}
