// Copyright (c) of parts are held by the various contributors (see the CLA)
// Licensed under the MIT License. See LICENSE file in the project root for full license information.
package polling

import (
	"io/ioutil"
	"net/http"
)

type OpenExchangeRates struct {
	Disclaimer string                 `json:"disclaimer"`
	License    string                 `json:"license"`
	Timestamp  string                 `json:"timestamp"`
	Base       string                 `json:"base"`
	Currency   OpenExchangeCurrencies `json:"rates"`
}

type OpenExchangeCurrencies struct {
	AED float64
	AFN float64
	ALL float64
	AMD float64
	ANG float64
	AOA float64
	ARS float64
	AUD float64
	AWG float64
	AZN float64
	BAM float64
	BBD float64
	BDT float64
	BGN float64
	BHD float64
	BIF float64
	BMD float64
	BND float64
	BOB float64
	BRL float64
	BSD float64
	BTC float64
	BTN float64
	BWP float64
	BYN float64
	BZD float64
	CAD float64
	CDF float64
	CHF float64
	CLF float64
	CLP float64
	CNH float64
	CNY float64
	COP float64
	CRC float64
	CUC float64
	CUP float64
	CVE float64
	CZK float64
	DJF float64
	DKK float64
	DOP float64
	DZD float64
	EGP float64
	ERN float64
	ETB float64
	EUR float64
	FJD float64
	FKP float64
	GBP float64
	GEL float64
	GGP float64
	GHS float64
	GIP float64
	GMD float64
	GNF float64
	GTQ float64
	GYD float64
	HKD float64
	HNL float64
	HRK float64
	HTG float64
	HUF float64
	IDR float64
	ILS float64
	IMP float64
	INR float64
	IQD float64
	IRR float64
	ISK float64
	JEP float64
	JMD float64
	JOD float64
	JPY float64
	KES float64
	KGS float64
	KHR float64
	KMF float64
	KPW float64
	KRW float64
	KWD float64
	KYD float64
	KZT float64
	LAK float64
	LBP float64
	LKR float64
	LRD float64
	LSL float64
	LYD float64
	MAD float64
	MDL float64
	MGA float64
	MKD float64
	MMK float64
	MNT float64
	MOP float64
	MRO float64
	MRU float64
	MUR float64
	MVR float64
	MWK float64
	MXN float64
	MYR float64
	MZN float64
	NAD float64
	NGN float64
	NIO float64
	NOK float64
	NPR float64
	NZD float64
	OMR float64
	PAB float64
	PEN float64
	PGK float64
	PHP float64
	PKR float64
	PLN float64
	PYG float64
	QAR float64
	RON float64
	RSD float64
	RUB float64
	RWF float64
	SAR float64
	SBD float64
	SCR float64
	SDG float64
	SEK float64
	SGD float64
	SHP float64
	SLL float64
	SOS float64
	SRD float64
	SSP float64
	STD float64
	SVC float64
	SYP float64
	SZL float64
	THB float64
	TJS float64
	TMT float64
	TND float64
	TOP float64
	TRY float64
	TTD float64
	TWD float64
	TZS float64
	UAH float64
	UGX float64
	USD float64
	UYU float64
	UZS float64
	VEF float64
	VES float64
	VND float64
	VUV float64
	WST float64
	XAF float64
	XAG float64
	XAU float64
	XCD float64
	XDR float64
	XOF float64
	XPD float64
	XPF float64
	XPT float64
	YER float64
	ZAR float64
	ZMW float64
	ZWL float64
}

//   you will need to replace the values put into peg structure
func CallOpenExchangeRates() ([]byte, error) {
	resp, err := http.Get("https://openexchangerates.org/api/latest.json?app_id=<INSERT API KEY HERE>")
	if err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return body, err
	}

}
