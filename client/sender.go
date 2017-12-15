package client

import (
	"encoding/xml"
	"log"
	"net"
)

//Sender : marshal node and send bytes to socket
func startSender(conn net.Conn, commandToSocket <-chan NCCCommand) {
	dataToSocket := make(chan []byte)

	go func() {
		for {
			data := <-dataToSocket
			conn.Write(data)
			conn.Write([]byte{delimeter})
			log.Println("Sent:\n" + string(data))
		}
	}()

	go func() {
		for {
			obj := <-commandToSocket
			data, err := xml.MarshalIndent(obj, "", "    ")
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			dataToSocket <- data
		}
	}()
}
