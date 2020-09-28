package converter

// The currency XML data
type currencyEnvelope struct {
	Sender string `xml:"Sender>name"`
	Cube   []cube `xml:"Cube>Cube>Cube"`
}

// The time XML data
type timeEnvelope struct {
	Time timeCube `xml:"Cube>Cube"`
}

// The time holder XML data
type timeCube struct {
	Time string `xml:"time,attr"`
}

// The cube XML structure
type cube struct {
	Name string  `xml:"currency,attr"`
	Rate float64 `xml:"rate,attr"`
}
