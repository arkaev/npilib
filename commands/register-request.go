package commands

import "encoding/xml"

type RegisterRq struct {
	XMLName xml.Name `xml:"NCC"`
	NCCCommand
	Request *registerRqRequest
}

type registerRqRequest struct {
	Name   string `xml:"name,attr"`
	Params *registerRqParams
}

type registerRqParams struct {
	ProtocolVersion int `xml:"protocol_version,attr"`
}

// CreateRegisterRequest will return "Register" command
func CreateRegisterRequest(protocolVersion int) NCCCommand {
	return &RegisterRq{
		Request: &registerRqRequest{
			Name:   "Register",
			Params: &registerRqParams{ProtocolVersion: protocolVersion}}}
}
