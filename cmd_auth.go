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
	NCCCommand
	Algorithm  string
	AuthScheme string
	Method     string
	Nonce      string
	Realm      string
	URI        string
	Username   string
}

type AuthenificateRqParser struct {
	Parser
}

//Unmarshal "Authenificate" command
func (h *AuthenificateRqParser) Unmarshal(data []byte) NCCCommand {
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

//AuthenificateHandler for "Authenificate" command
type AuthenificateHandler struct {
	Handler
	conn *Conn
}

//Handle "Authenificate" command
func (h *AuthenificateHandler) Handle(cmd NCCCommand) {
	auth := cmd.(*AuthenificateRq)

	value := calculateMD5(h.conn.digest, auth.Nonce, auth.Method, auth.URI)
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
