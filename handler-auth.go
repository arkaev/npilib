package npilib

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"

	c "github.com/arkaev/npilib/commands"
)

//HandleAuthenificate will process "Authenificate" message
func HandleAuthenificate(nc *Conn, msg *Msg) {
	auth := msg.Parsed.(*c.AuthentificateRq)
	params := auth.Request.Params

	value := calculateMD5(nc.digest, params.Nonce, params.Method, params.URI)
	value = strings.ToLower(value)

	nc.Publish(c.CreateAuthenticateResponse(value))
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
