package commands

import "encoding/xml"

type RegisterRq struct {
	XMLName xml.Name `xml:"NCC"`
	NCCCommand
	Request *RegisterRqRequest
}

type RegisterRqRequest struct {
	Name   string `xml:"name,attr"`
	Params *RegisterRqParams
}

type RegisterRqParams struct {
	ProtocolVersion int `xml:"protocol_version,attr"`
}
