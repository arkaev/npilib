package npilib

import c "github.com/arkaev/npilib/commands"

//EchoHandler for "Echo" command
type EchoHandler struct {
	Handler
	conn *Conn
}

//Handle "Echo" command
func (h *EchoHandler) Handle(cmd c.NCCCommand) {
	nccn := &c.EchoRs{Response: &c.EchoRsResponse{Name: "Echo"}}

	h.conn.commandToSocket <- nccn
}
