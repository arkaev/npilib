package npilib

import (
	"encoding/xml"
)

type RegisterPeerRq struct {
	NCCCommand
	XMLName xml.Name `xml:"NCCN"`
	Request *RegisterPeerRqResponse
}

type RegisterPeerRqResponse struct {
	Params *RegisterPeerRqParams
}

type RegisterPeerRqParams struct {
	Login       string `xml:"login,attr"`
	MaxProtocol int    `xml:"max_protocol,attr"`
	MinProtocol string `xml:"min_protocol,attr"`
	Role        string `xml:"role,attr"`
}

type RegisterPeerRs struct {
	NCCCommand
	XMLName  xml.Name `xml:"NCCN"`
	Response *RegisterPeerRsResponse
}

type RegisterPeerRsResponse struct {
	Params *RegisterPeerRsParams
}

type RegisterPeerRsParams struct {
	AllowEncoding   string `xml:"allow_encoding,attr"`
	Domain          string `xml:"domain,attr"`
	Node            string `xml:"node,attr"`
	Peer            string `xml:"peer,attr"`
	ProtocolVersion int    `xml:"protocol_version,attr"`
}

type RegisterPeerRsParser struct {
	Parser
}

func (p *RegisterPeerRsParser) Unmarshal(data []byte) NCCCommand {
	var rq RegisterPeerRs
	xml.Unmarshal(data, &rq)

	return &rq
}

//RegisterPeerHandler for "RegisterPeer" command
type RegisterPeerHandler struct {
	Handler
	conn *Conn
}

//Handle "RegisterPeer" command
func (h *RegisterPeerHandler) Handle(cmd NCCCommand) {

	rs := cmd.(*RegisterPeerRs).Response.Params
	h.conn.allowEncoding = rs.AllowEncoding
	h.conn.domain = rs.Domain
	h.conn.node = rs.Node
	h.conn.peer = rs.Peer
	h.conn.protocolVersion = rs.ProtocolVersion

	type Params struct {
		ProtocolVersion int `xml:"protocol_version,attr"`
	}

	type Request struct {
		Name   string `xml:"name,attr"`
		Params *Params
	}

	type NCC struct {
		NCCCommand
		Request *Request
	}

	h.conn.commandToSocket <- &NCC{
		Request: &Request{
			Name:   "Register",
			Params: &Params{ProtocolVersion: 600}}}
}
