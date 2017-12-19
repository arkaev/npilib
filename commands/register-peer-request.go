package commands

import (
	"encoding/xml"
)

type RegisterPeerRq struct {
	NCCCommand
	XMLName xml.Name `xml:"NCCN"`
	Request *RegisterPeerRqRequest
}

type RegisterPeerRqRequest struct {
	Name   string `xml:"name,attr"`
	Params *RegisterPeerRqParams
}

type RegisterPeerRqParams struct {
	Login       string `xml:"login,attr"`
	MaxProtocol int    `xml:"max_protocol,attr"`
	MinProtocol int    `xml:"min_protocol,attr"`
	Role        string `xml:"role,attr"`
}
