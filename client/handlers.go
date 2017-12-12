package client

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

//Handler process command with name
type Handler interface {
	Process(node *Node)
}

//AuthenificateHandler for "Authenificate" command
type AuthenificateHandler struct {
	Handler
	digest      string
	outCommands chan<- NCCCommand
}

//Process "Authenificate" command
func (h *AuthenificateHandler) Process(node *Node) {
	//algoritm := node.Nodes[0].Attributes["algoritm"]
	//auth_scheme := node.Nodes[0].Attributes["auth_scheme"]
	method := node.Nodes[0].Attributes["method"]
	nonce := node.Nodes[0].Attributes["nonce"]
	//realm := node.Nodes[0].Attributes["realm"]
	uri := node.Nodes[0].Attributes["uri"]
	//username := node.Nodes[0].Attributes["username"]

	value := calculateMD5(h.digest, nonce, method, uri)
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

	h.outCommands <- &NCCN{
		Response: &Response{
			Name:  "Authenticate",
			Param: &Param{Name: "response", Value: value}}}
}

//RegisterPeerHandler for "RegisterPeer" command
type RegisterPeerHandler struct {
	Handler
	outCommands chan<- NCCCommand
}

//Process "RegisterPeer" command
func (h *RegisterPeerHandler) Process(node *Node) {
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

	h.outCommands <- &NCC{
		Request: &Request{
			Name:   "Register",
			Params: &Params{ProtocolVersion: 600}}}
}

//EchoHandler for "Echo" command
type EchoHandler struct {
	Handler
	outCommands chan<- NCCCommand
}

//Process "Echo" command
func (h *EchoHandler) Process(node *Node) {
	type Response struct {
		Name string `xml:"name,attr"`
	}

	type NCCN struct {
		NCCCommand
		Response *Response
	}

	nccn := &NCCN{Response: &Response{Name: "Echo"}}

	h.outCommands <- nccn
}

//RegisterHandler for "Register" command
type RegisterHandler struct {
	Handler
}

//Process "Register" command
func (h *RegisterHandler) Process(node *Node) {
	fmt.Println("Successful registration")
	//TODO add handlers
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
