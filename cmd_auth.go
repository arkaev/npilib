package npilib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
)

//AuthenificateHandler for "Authenificate" command
type AuthenificateHandler struct {
	Handler
	conn *Conn

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
	value := calculateMD5(h.conn.digest, h.nonce, h.method, h.uri)
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

	h.conn.commandToSocket <- &NCCN{
		Response: &Response{
			Name:  "Authenticate",
			Param: &Param{Name: "response", Value: value}}}
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
