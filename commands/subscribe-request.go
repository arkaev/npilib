package commands

import "encoding/xml"

type SubscribeRq struct {
	XMLName xml.Name `xml:"NCC"`
	Request *SubscribeRqRequest
}

type SubscribeRqRequest struct {
	Name   string `xml:"name,attr"`
	Params []*SubscribeRqParams
}

type SubscribeRqParams struct {
	List    string `xml:"list,attr"`
	Enabled bool   `xml:"enabled,attr"`
	Instant bool   `xml:"instant,attr"`
}
