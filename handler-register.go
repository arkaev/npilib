package npilib

import (
	"log"

	c "github.com/arkaev/npilib/commands"
)

//RegisterHandler for "Register" command
type RegisterHandler struct {
	Handler
	conn *Conn
}

//Handle "Register" command
func (h *RegisterHandler) Handle(cmd c.NCCCommand) {
	log.Println("Successful registration")

	h.conn.commandToSocket <- SubscribeCommand("callslist")
	h.conn.commandToSocket <- SubscribeCommand("buddylist")
}
