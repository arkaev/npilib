package npilib

import "log"

//RegisterHandler for "Register" command
type RegisterHandler struct {
	Handler
	conn *Conn
}

//Handle "Register" command
func (h *RegisterHandler) Handle(cmd NCCCommand) {
	log.Println("Successful registration")

	h.conn.commandToSocket <- SubscribeCommand("callslist")
	h.conn.commandToSocket <- SubscribeCommand("buddylist")
}
