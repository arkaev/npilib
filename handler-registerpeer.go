package npilib

import (
	"encoding/xml"

	c "github.com/arkaev/npilib/commands"
)

type RegisterPeerRsParser struct {
	Parser
}

func (p *RegisterPeerRsParser) Unmarshal(data []byte) c.NCCCommand {
	var rq c.RegisterPeerRs
	xml.Unmarshal(data, &rq)

	return &rq
}

//RegisterPeerHandler for "RegisterPeer" command
type RegisterPeerHandler struct {
	Handler
	conn *Conn
}

//Handle "RegisterPeer" command
func (h *RegisterPeerHandler) Handle(cmd c.NCCCommand) {

	rs := cmd.(*c.RegisterPeerRs).Response.Params
	h.conn.allowEncoding = rs.AllowEncoding
	h.conn.domain = rs.Domain
	h.conn.node = rs.Node
	h.conn.peer = rs.Peer
	h.conn.protocolVersion = rs.ProtocolVersion

	h.conn.commandToSocket <- &c.RegisterRq{
		Request: &c.RegisterRqRequest{
			Name:   "Register",
			Params: &c.RegisterRqParams{ProtocolVersion: 600}}}
}
