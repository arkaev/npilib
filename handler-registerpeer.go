package npilib

import (
	c "github.com/arkaev/npilib/commands"
)

//HandleRegisterPeer will process "RegisterPeer" command
func HandleRegisterPeer(nc *Conn, msg *Msg) {
	rs := msg.Parsed.(*c.RegisterPeerRs).Response.Params
	nc.allowEncoding = rs.AllowEncoding
	nc.domain = rs.Domain
	nc.node = rs.Node
	nc.peer = rs.Peer
	nc.protocolVersion = rs.ProtocolVersion

	nc.Publish(c.CreateRegisterRequest(600))
}
