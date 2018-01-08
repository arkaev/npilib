package npilib

import (
	"encoding/xml"
	"log"

	c "github.com/arkaev/npilib/commands"
)

//Sender : marshal node and send bytes to socket
func startSender(nc *Conn) {
	dataToSocket := make(chan []byte)

	go func(in chan []byte, client *Conn) {
		for data := range in {
			client.send(data)
			log.Printf("Sent:\n%s\n", data)
		}
	}(dataToSocket, nc)

	go func(in chan c.NCCCommand, out chan []byte) {
		for obj := range in {
			data, err := xml.MarshalIndent(obj, "", "    ")
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			out <- data
		}
	}(nc.commandToSocket, dataToSocket)
}
