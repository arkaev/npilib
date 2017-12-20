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

//CreateRegisterPeerCommand will construct RegisterPeer command
func CreateRegisterPeerCommand(login string) NCCCommand {
	return &RegisterPeerRq{
		Request: &RegisterPeerRqRequest{Name: "RegisterPeer",
			Params: &RegisterPeerRqParams{
				Login:       login,
				MaxProtocol: 0,
				MinProtocol: 0,
				Role:        "service"}}}
}
