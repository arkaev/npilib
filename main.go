package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"net"
	"time"
)

const delimeter byte = 0
const address string = "docker72:3242"

type NCCN struct {
	Request Request
}

type Request struct {
	Name   string `xml:"name,attr"`
	Params Params
}

type Params struct {
	Login       string `xml:"login,attr"`
	MaxProtocol int    `xml:"max_protocol,attr"`
	MinProtocol int    `xml:"min_protocol,attr"`
	Role        string `xml:"role,attr"`
}

func main() {
	fmt.Println("Started")

	conn, err := connect(address)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	time.Sleep(time.Second * 3)

	fmt.Println("Exit")
}

func sender(conn net.Conn, commands <-chan NCCN) {
	cmd := <-commands
	output, err := xml.MarshalIndent(cmd, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	conn.Write(output)
	conn.Write([]byte{delimeter})
	fmt.Print("Sent: ")
	fmt.Println(string(output))
}

func receiver(conn net.Conn) {
	var inQueue = make(chan string)
	go commandParser(inQueue)

	bufReader := bufio.NewReader(conn)
	for {
		msg, err := bufReader.ReadString(delimeter)
		if err != nil {
			break
		}
		fmt.Print("Received: ")
		fmt.Println(msg)
		inQueue <- msg
	}
}

func commandParser(strCommands <-chan string) {
	// msg := <-strCommands
}

func connect(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected")

	var outChannel = make(chan NCCN)

	go sender(conn, outChannel)
	go receiver(conn)

	params := Params{Login: "naucrm", MaxProtocol: 0, MinProtocol: 0, Role: "service"}
	registerPeer := Request{Name: "RegisterPeer", Params: params}
	registerPeerNccn := NCCN{Request: registerPeer}

	outChannel <- registerPeerNccn

	return conn, nil
}
