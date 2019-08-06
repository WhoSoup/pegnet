package common_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	. "github.com/pegnet/pegnet/common"
)

func TestAmounts(t *testing.T) {
	vectors := []struct {
		V int64
		S string
	}{ // TODO: Add vectors
		{1e8, "1"},
		{2e8, "2"},
		{2e8 + 2e7, "2.2"},
		{1, "0.00000001"},
		{0, "0"},
		{12345678, "0.12345678"},
	}

	for _, v := range vectors {
		// Test the expect
		vS := AmountToString(v.V)
		if vS != v.S {
			t.Errorf("[1] Exp %s, got %s", v.S, vS)
		}

		vV, err := StringToAmount(v.S)
		if err != nil {
			t.Errorf("[2] Exp %d, got %d", v.V, vV)
		}

		// Test the results
		if vS2 := AmountToString(vV); vS2 != vS {
			t.Errorf("[3] Exp %s, got %s", vS2, vS)
		}

		if vV2, _ := StringToAmount(vS); vV2 != vV {
			t.Errorf("[3] Exp %d, got %d", vV2, vV)
		}
	}
}

func TestAmountJsonMarshal(t *testing.T) {
	type TestStruct struct {
		Amt Amount
	}

	for i := uint64(0); i < 100000; i++ {
		ts := &TestStruct{Amount(rand.Int63())}
		d, err := json.Marshal(ts)
		if err != nil {
			t.Error(err)
		}

		t2 := new(TestStruct)
		err = json.Unmarshal(d, t2)
		if err != nil {
			t.Error(err)
		}

		if ts.Amt != t2.Amt {
			fmt.Println(string(d))
			t.Error("json failed")
			t.FailNow()
		}
	}

}

func TestFromFloat(t *testing.T) {
	for i := 0; i < 1000; i++ {
		// test the string
		f := rand.Float64()
		v := FloatToAmount(f)

		// truncate f, so it does not round
		f = math.Trunc(f*1e8) / 1e8

		fS := fmt.Sprintf("%.8f", f)
		fS = strings.TrimRight(fS, "0")
		vS := AmountToString(v)

		if fS != vS {
			t.Errorf("Exp %s, got %s", fS, vS)
		}
	}
}

func TestParseFixed(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", true},
		{"zero", args{"0"}, "0.00000000", false},
		{"one", args{"1"}, "1.00000000", false},
		{"negative", args{"-1"}, "", true},
		{"one decimal", args{".1"}, "0.10000000", false},
		{"all decimal", args{"0.12345678"}, "0.12345678", false},
		{"decimal no pre0", args{".12345678"}, "0.12345678", false},
		{"too many decimal", args{"0.123456789"}, "", true},
		{"large non frac", args{"12303580123891845729385719382479835723984723948723948723498723498273450928345234"}, "12303580123891845729385719382479835723984723948723948723498723498273450928345234.00000555", false},
		{"large frac", args{"12303580123891845729385719382479835723984723948723948723498723498273450928345234"}, "12303580123891845729385719382479835723984723948723948723498723498273450928345234.00000555", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFixed(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFixed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.FloatString(8) != tt.want {
				t.Errorf("ParseFixed() = %v, want %v", got.FloatString(8), tt.want)
			}
		})
	}
}

func TestTrimFixed(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"zero", "0", "0"},
		{"hundred", "100", "100"},
		{"decimal", "0.1", "0.1"},
		{"trim all decimal", "1.000", "1"},
		{"trim some decimal", "1.020", "1.02"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := ParseFixed(tt.args)
			if got := TrimFixed(r); got != tt.want {
				t.Errorf("TrimFixed() = %v, want %v", got, tt.want)
			}
		})
	}
}
