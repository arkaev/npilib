package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io"
	"strings"
)

//TODO use key file
const digestTemp = "787013196de72b3567ff0f3aac0394a6"

func HandleAuthenificate(node *Node, outCommands chan<- *Node) {
	//algoritm := node.Nodes[0].Attributes["algoritm"]
	//auth_scheme := node.Nodes[0].Attributes["auth_scheme"]
	method := node.Nodes[0].Attributes["method"]
	nonce := node.Nodes[0].Attributes["nonce"]
	//realm := node.Nodes[0].Attributes["realm"]
	uri := node.Nodes[0].Attributes["uri"]
	//username := node.Nodes[0].Attributes["username"]
	digest := digestTemp

	value := calculateMD5(digest, nonce, method, uri)
	value = strings.ToLower(value)

	paramAttrs := make(map[string]string)
	paramAttrs["name"] = "response"
	paramAttrs["value"] = value
	paramNode := &Node{
		XMLName:    xml.Name{Local: "Param"},
		Attributes: paramAttrs}

	requestAttrs := make(map[string]string)
	requestAttrs["name"] = "Authenticate"
	requestNode := &Node{
		XMLName:    xml.Name{Local: "Response"},
		Attributes: requestAttrs,
		Nodes:      []*Node{paramNode}}

	rootNode := &Node{
		XMLName: xml.Name{Local: "NCCN"},
		Nodes:   []*Node{requestNode}}

	outCommands <- rootNode
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
