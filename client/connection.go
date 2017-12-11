package client

import (
	"encoding/xml"
	"fmt"
	"net"
)

const delimeter byte = 0

//Connect: create connection by address and keyFile
func Connect(address string, keyFile string) (net.Conn, error) {

	md5, err := GetDigest(keyFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Digest: " + md5)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected")

	var outChannel = make(chan *Node)

	go Sender(conn, outChannel)
	go Receiver(conn, outChannel)

	outChannel <- registerPeer()

	return conn, nil
}

func registerPeer() *Node {
	paramAttrs := make(map[string]string)
	paramAttrs["login"] = "naucrm"
	paramAttrs["max_protocol"] = "0"
	paramAttrs["min_protocol"] = "0"
	paramAttrs["role"] = "service"
	param := &Node{
		XMLName:    xml.Name{Local: "Params"},
		Attributes: paramAttrs}

	requestAttrs := make(map[string]string)
	requestAttrs["name"] = "RegisterPeer"
	request := &Node{
		XMLName:    xml.Name{Local: "Request"},
		Nodes:      []*Node{param},
		Attributes: requestAttrs}

	return &Node{
		XMLName: xml.Name{Local: "NCCN"},
		Nodes:   []*Node{request}}
}
