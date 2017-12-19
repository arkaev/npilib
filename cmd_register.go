package npilib

import "log"

//RegisterHandler for "Register" command
type RegisterHandler struct {
	Handler
	conn *Conn
}

//Unmarshal "Register" command
func (h *RegisterHandler) Unmarshal(node *Node) Handler {
	return h
}

//Handle "Register" command
func (h *RegisterHandler) Handle() {
	log.Println("Successful registration")

	h.conn.commandToSocket <- SubscribeCommand("callslist")
	h.conn.commandToSocket <- SubscribeCommand("buddylist")
}
