package npilib

import (
	"encoding/xml"
	"log"
)

//Sender : marshal node and send bytes to socket
func startSender(nc *Conn) {
	dataToSocket := make(chan []byte)

	go func() {
		for {
			data := <-dataToSocket
			nc.send(data)
			log.Println("Sent:\n" + string(data))
		}
	}()

	go func() {
		for {
			obj := <-nc.commandToSocket
			data, err := xml.MarshalIndent(obj, "", "    ")
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			dataToSocket <- data
		}
	}()
}
