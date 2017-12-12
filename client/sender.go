package client

import (
	"encoding/xml"
	"fmt"
	"net"
)

//Sender : marshal node and send bytes to socket
func Sender(conn net.Conn, objectToSocket <-chan NCCCommand) {
	dataToSocket := make(chan []byte)
	go sendBytesToSocket(conn, dataToSocket)

	for {
		obj := <-objectToSocket
		data, err := xml.MarshalIndent(obj, "", "    ")
		if err != nil {
			fmt.Println("error: %v", err)
		}
		dataToSocket <- data
	}
}

func sendBytesToSocket(conn net.Conn, dataToSocket <-chan []byte) {
	for {
		data := <-dataToSocket
		conn.Write(data)
		conn.Write([]byte{delimeter})
		fmt.Println("Sent:\n" + string(data))
	}
}
