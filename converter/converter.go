package converter

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"
)

const(
	currencyDateFormat = "2006-01-02" // The time format in ECB XML
	eur                = "EUR"        // The euro symbol
)

func FetchCurrencyData() (data []byte, err error) {
	res, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close() // ignore error?
	return data, err
}

func ParseCurrencyData(data []byte) (ts time.Time, currencies map[string]float64, err error) {
	// parse once to get the currencies, return on error
	var c currencyEnvelope
	err = xml.Unmarshal(data, &c)
	if err != nil {
		return time.Time{}, nil, err
	}

	// parse again to get the timestamp, return on error
	var t timeEnvelope
	err = xml.Unmarshal(data, &t)
	if err != nil {
		return time.Time{}, nil, err
	}

	// parse time, return on error
	ts, err = time.Parse(currencyDateFormat, t.Time.Time)
	if err != nil {
		return time.Time{}, nil, err
	}

	currencies = make(map[string]float64)

	// manually insert EUR as "1"
	currencies[eur] = 1

	// insert all rates
	for _, currency := range c.Cube {
		currencies[currency.Name] = currency.Rate
	}

	return ts, currencies, nil
}
