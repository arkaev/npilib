package npilib

//EchoHandler for "Echo" command
type EchoHandler struct {
	Handler
	conn *Conn
}

//Unmarshal "Echo" command
func (h *EchoHandler) Unmarshal(node *Node) Handler {
	return h
}

//Handle "Echo" command
func (h *EchoHandler) Handle() {
	type Response struct {
		Name string `xml:"name,attr"`
	}

	type NCCN struct {
		NCCCommand
		Response *Response
	}

	nccn := &NCCN{Response: &Response{Name: "Echo"}}

	h.conn.commandToSocket <- nccn
}
