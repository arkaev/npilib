package npilib

import (
	"log"
	"strconv"
)

//Handler process command with name
type Handler interface {
	//Unmarshal node to command pojo
	Unmarshal(node *Node) Handler
	//Handle command and process
	Handle()
}

//CommonTagHandler if handlers have same root tag but different name-attribute
type CommonTagHandler struct {
	Handler
	handlers map[string]Handler
}

//Unmarshal command for wrapper
func (h *CommonTagHandler) Unmarshal(event *Node) Handler {
	name := event.Attributes["name"]
	handler, exist := h.handlers[name]
	if exist {
		wrapped := handler.Unmarshal(event)
		return wrapped
	}

	log.Printf("Unknown handler: %s\n", name)
	return nil
}

//Handle command for wrapper
func (h *CommonTagHandler) Handle() {
}

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

//DoNothingHandler bulk handler
type DoNothingHandler struct {
	Handler
}

//Unmarshal bulk
func (h *DoNothingHandler) Unmarshal(node *Node) Handler {
	return h
}

//Handle bulk
func (h *DoNothingHandler) Handle() {
}
