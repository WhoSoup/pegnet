// Copyright (c) of parts are held by the various contributors (see the CLA)
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
package opr_test

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/FactomProject/btcutil/base58"
	. "github.com/pegnet/pegnet/opr"
)

func TestOPR_JSON_Marshal(t *testing.T) {
	LX.Init(0x123412341234, 25, 256, 5)
	opr := new(OraclePriceRecord)

	opr.Difficulty = 1
	opr.Grade = 1
	//opr.Nonce = base58.Encode(LX.Hash([]byte("a Nonce")))
	//opr.ChainID = base58.Encode(LX.Hash([]byte("a chainID")))
	opr.Dbht = 1901232
	opr.WinPreviousOPR = [10]string{
		base58.Encode(LX.Hash([]byte("winner number 1"))),
		base58.Encode(LX.Hash([]byte("winner number 2"))),
		base58.Encode(LX.Hash([]byte("winner number 3"))),
		base58.Encode(LX.Hash([]byte("winner number 4"))),
		base58.Encode(LX.Hash([]byte("winner number 5"))),
		base58.Encode(LX.Hash([]byte("winner number 6"))),
		base58.Encode(LX.Hash([]byte("winner number 7"))),
		base58.Encode(LX.Hash([]byte("winner number 8"))),
		base58.Encode(LX.Hash([]byte("winner number 9"))),
		base58.Encode(LX.Hash([]byte("winner number 10"))),
	}
	opr.CoinbasePNTAddress = "pPNT4wBqpZM9xaShSYTABzAf1i1eSHVbbNk2xd1x6AkfZiy366c620f"
	opr.FactomDigitalID = []string{"miner", "one"}
	opr.PNT = 2
	opr.USD = 20
	opr.EUR = 200
	opr.JPY = 11
	opr.GBP = 12
	opr.CAD = 13
	opr.CHF = 14
	opr.INR = 15
	opr.SGD = 16
	opr.CNY = 17
	opr.HKD = 18
	opr.XAU = 19
	opr.XAG = 101
	opr.XPD = 1012
	opr.XPT = 10123
	opr.XBT = 10124
	opr.ETH = 10125
	opr.LTC = 10126
	opr.XBC = 10127
	opr.FCT = 10128

	v, _ := json.Marshal(opr)
	fmt.Println("len of entry", len(string(v)), "\n\n", string(v))
	opr2 := new(OraclePriceRecord)
	json.Unmarshal(v, &opr2)
	v2, _ := json.Marshal(opr2)
	fmt.Println("\n\n", string(v2))
	if string(v2) != string(v) {
		t.Error("JSON is different")
	}
}

func dummyCall(b []byte) {

}

func len2(b int) int {
	if b&0xFF000000 != 0 {
		return 0
	}
	if b&0xFFFF0000 != 0 {
		return 1
	}
	if b&0xFFFFFF00 != 0 {
		return 2
	}
	return 3
}
func BenchmarkNonce(b *testing.B) {

	b.Run("old nonce", func(c *testing.B) {
		nonce := []byte{0, 0}
		for i := 0; i < c.N; i++ {
			nonce = nonce[:0]
			for j := i; j > 0; j = j >> 8 {
				nonce = append(nonce, byte(j))
			}
		}
		dummyCall(nonce)
	})
	b.Run("new nonce A", func(c *testing.B) {
		nonce := []byte{0}
		pos := 0
		for i := 0; i < c.N; i++ {
			nonce[pos]++
			for j := pos; j >= 0; j-- {
				if nonce[pos] == 0 {
					if pos > 0 {
						nonce[pos-1]++
					} else {
						nonce = append([]byte{1}, nonce...)
						pos++
					}
				} else {
					break
				}
			}
			dummyCall(nonce)
		}
	})

	b.Run("new nonce B", func(c *testing.B) {
		nonce := make([]byte, 4)
		for i := 0; i < c.N; i++ {
			binary.BigEndian.PutUint32(nonce, uint32(i))
			first := len2(i)
			dummyCall(nonce[first:])
		}

	})
}
