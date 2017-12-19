package npilib

import "strconv"

//RegisterPeerHandler for "RegisterPeer" command
type RegisterPeerHandler struct {
	Handler
	conn            *Conn
	AllowEncoding   string
	Domain          string
	Node            string
	Peer            string
	ProtocolVersion int
}

//Unmarshal "RegisterPeer" command
func (h *RegisterPeerHandler) Unmarshal(node *Node) Handler {
	paramsNode := node.Nodes[0]

	h.AllowEncoding = paramsNode.Attributes["allow_encoding"]
	h.Domain = paramsNode.Attributes["domain"]
	h.Node = paramsNode.Attributes["node"]
	h.Peer = paramsNode.Attributes["peer"]
	h.ProtocolVersion, _ = strconv.Atoi(paramsNode.Attributes["protocol_version"])

	return h
}

//Handle "RegisterPeer" command
func (h *RegisterPeerHandler) Handle() {
	h.conn.allowEncoding = h.AllowEncoding
	h.conn.domain = h.Domain
	h.conn.node = h.Node
	h.conn.peer = h.Peer
	h.conn.protocolVersion = h.ProtocolVersion

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
