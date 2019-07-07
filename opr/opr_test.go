// Copyright (c) of parts are held by the various contributors (see the CLA)
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
package opr_test

import (
	"crypto/rand"
	"encoding/hex"
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

func BenchmarkRefactoredFromSlice(b *testing.B) {
	var suma, sumb uint64
	b.Run("slice to uint64 (new)", func(b *testing.B) {
		var h [64]byte
		b.StopTimer()
		rand.Read(h[:])
		b.StartTimer()
		for xx := 0; xx < b.N; xx++ {
			sumb = uint64(h[7]) | uint64(h[6])<<8 | uint64(h[5])<<16 | uint64(h[4])<<24 |
				uint64(h[3])<<32 | uint64(h[2])<<40 | uint64(h[1])<<48 | uint64(h[0])<<56
		}
	})

	b.Run("slice to uint64 (old)", func(b *testing.B) {
		var h [64]byte
		b.StopTimer()
		rand.Read(h[:])
		b.StartTimer()
		for xx := 0; xx < b.N; xx++ {
			suma = 0
			for i := uint64(0); i < 8; i++ {
				suma = suma<<8 + uint64(h[i])
			}
		}
	})

	_ = sumb
}

func BenchmarkRefactoredToSlice(b *testing.B) {
	b.Run("int to slice (new)", func(b *testing.B) {
		nonce := make([]byte, 4)
		for xx := 0; xx < b.N; xx++ {
			nonce[0] = byte(xx >> 24)
			nonce[1] = byte(xx >> 16)
			nonce[2] = byte(xx >> 8)
			nonce[3] = byte(xx)
		}
	})
	b.Run("int to slice (old)", func(b *testing.B) {
		nonce := []byte{0, 0}
		for xx := 0; xx < b.N; xx++ {
			nonce = nonce[:0]
			for j := b.N; j > 0; j = j >> 8 {
				nonce = append(nonce, byte(j))
			}
		}
	})
}

func TestOraclePriceRecord_ComputeDifficulty_Change(t *testing.T) {
	var h []byte
	h = make([]byte, 8)
	var suma, sumb uint64
	for i := 0; i < 50; i++ {
		rand.Read(h)
		suma = 0
		for i := uint64(0); i < 8; i++ {
			suma = suma<<8 + uint64(h[i])
		}
		sumb = uint64(h[7]) | uint64(h[6])<<8 | uint64(h[5])<<16 | uint64(h[4])<<24 |
			uint64(h[3])<<32 | uint64(h[2])<<40 | uint64(h[1])<<48 | uint64(h[0])<<56
		if suma != sumb {
			t.Errorf("two results not the same. input = %s, old = %x, new = %x", hex.EncodeToString(h), suma, sumb)
		}
	}
}
