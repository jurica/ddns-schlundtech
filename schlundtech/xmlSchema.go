package schlundtech

import "encoding/xml"

type Response struct {
	XMLName xml.Name `xml:"response"`
	Result  Result   `xml:"result"`
}

type Result struct {
	XMLName xml.Name `xml:"result"`
	Data    Data     `xml:"data"`
}

type Data struct {
	XMLName xml.Name `xml:"data"`
	Zone    Zone     `xml:"zone"`
}

type Zone struct {
	XMLName  xml.Name `xml:"zone"`
	Name     string   `xml:"name"`
	Rrs      []RR     `xml:"rr"`
	SystemNS string   `xml:"system_ns"`
}

type RR struct {
	XMLName xml.Name `xml:"rr"`
	Name    string   `xml:"name"`
	Ttl     string   `xml:"ttl"`
	Type    string   `xml:"type"`
	Value   string   `xml:"value"`
}
