package npilib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io"
	"strings"
)

// AuthenificateRq structure contains "Authenificate" request data
type AuthenificateRq struct {
	Algorithm  string
	AuthScheme string
	Method     string
	Nonce      string
	Realm      string
	URI        string
	Username   string
}

//AuthenificateHandler for "Authenificate" command
type AuthenificateHandler struct {
	Handler
	conn *Conn

	algorithm  string
	authScheme string
	method     string
	nonce      string
	realm      string
	uri        string
	username   string
}

//Unmarshal "Authenificate" command
func (h *AuthenificateHandler) Unmarshal(node *Node) Handler {
	h.algorithm = node.Nodes[0].Attributes["algoritm"]
	h.authScheme = node.Nodes[0].Attributes["auth_scheme"]
	h.method = node.Nodes[0].Attributes["method"]
	h.nonce = node.Nodes[0].Attributes["nonce"]
	h.realm = node.Nodes[0].Attributes["realm"]
	h.uri = node.Nodes[0].Attributes["uri"]
	h.username = node.Nodes[0].Attributes["username"]

	return h
}

//Parse "Authenificate" command
func (h *AuthenificateHandler) Parse(data []byte) *AuthenificateRq {
	type Params struct {
		Algorithm  string `xml:"algorithm,attr"`
		AuthScheme string `xml:"auth_scheme,attr"`
		Method     string `xml:"method,attr"`
		Nonce      string `xml:"nonce,attr"`
		Realm      string `xml:"realm,attr"`
		URI        string `xml:"uri,attr"`
		Username   string `xml:"username,attr"`
	}

	type Request struct {
		Name   string `xml:"name,attr"`
		Params *Params
	}

	type NCCN struct {
		Request *Request
	}

	var auth NCCN
	xml.Unmarshal(data, &auth)

	params := auth.Request.Params

	return &AuthenificateRq{
		Algorithm:  params.Algorithm,
		AuthScheme: params.AuthScheme,
		Method:     params.Method,
		Nonce:      params.Nonce,
		Realm:      params.Realm,
		URI:        params.URI,
		Username:   params.Username,
	}
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
