package npilib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"strings"
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

//AuthenificateHandler for "Authenificate" command
type AuthenificateHandler struct {
	Handler
	digest string

	out        chan<- NCCCommand
	algoritm   string
	authScheme string
	method     string
	nonce      string
	realm      string
	uri        string
	username   string
}

//Unmarshal "Authenificate" command
func (h *AuthenificateHandler) Unmarshal(node *Node) Handler {
	h.algoritm = node.Nodes[0].Attributes["algoritm"]
	h.authScheme = node.Nodes[0].Attributes["auth_scheme"]
	h.method = node.Nodes[0].Attributes["method"]
	h.nonce = node.Nodes[0].Attributes["nonce"]
	h.realm = node.Nodes[0].Attributes["realm"]
	h.uri = node.Nodes[0].Attributes["uri"]
	h.username = node.Nodes[0].Attributes["username"]

	return h
}

//Handle "Authenificate" command
func (h *AuthenificateHandler) Handle() {
	value := calculateMD5(h.digest, h.nonce, h.method, h.uri)
	value = strings.ToLower(value)

	type Param struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}

	type Response struct {
		Name  string `xml:"name,attr"`
		Param *Param
	}

	type NCCN struct {
		NCCCommand
		Response *Response
	}

	h.out <- &NCCN{
		Response: &Response{
			Name:  "Authenticate",
			Param: &Param{Name: "response", Value: value}}}
}

//RegisterPeerHandler for "RegisterPeer" command
type RegisterPeerHandler struct {
	Handler
	config          *RegistrationInfo
	out             chan<- NCCCommand
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
	h.config.AllowEncoding = h.AllowEncoding
	h.config.Domain = h.Domain
	h.config.Node = h.Node
	h.config.Peer = h.Peer
	h.config.ProtocolVersion = h.ProtocolVersion

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

	h.out <- &NCC{
		Request: &Request{
			Name:   "Register",
			Params: &Params{ProtocolVersion: 600}}}
}

//EchoHandler for "Echo" command
type EchoHandler struct {
	Handler
	out chan<- NCCCommand
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

	h.out <- nccn
}

//RegisterHandler for "Register" command
type RegisterHandler struct {
	Handler
	out chan<- NCCCommand
}

//Unmarshal "Register" command
func (h *RegisterHandler) Unmarshal(node *Node) Handler {
	return h
}

//Handle "Register" command
func (h *RegisterHandler) Handle() {
	log.Println("Successful registration")

	h.out <- SubscribeCommand("callslist")
	h.out <- SubscribeCommand("buddylist")
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

func calculateMD5(digest, nonce, method, uri string) string {
	return hashToString(digest + ":" + nonce + ":" + hashToString(method+":"+uri))
}

func hashToString(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
