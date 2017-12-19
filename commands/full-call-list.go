package commands

import "encoding/xml"

type FullCallsList struct {
	XMLName       xml.Name `xml:"NCC"`
	From          string   `xml:"from,attr"`
	To            string   `xml:"to,attr"`
	FullCallsList *FullCallsListMain
}

type FullCallsListMain struct {
	TimeT uint64 `xml:"time_t,attr"`
}
