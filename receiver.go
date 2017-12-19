package npilib

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"time"
)

//Receiver for commands from socket
func startReceiver(nc *Conn) {
	socketToDataCommand := make(chan []byte)
	dataToNode := make(chan *Node)
	nodeToHanlderChannel := make(chan Handler)

	go func() {
		bufReader := bufio.NewReader(nc.conn)
		for {
			cmdData, err := bufReader.ReadBytes(delimeter)
			if err != nil {
				if err == io.EOF {
					//sleep if no data
					time.Sleep(time.Millisecond * 10)
				} else {
					log.Printf("Unexpected read error: %s\n", err)
					break
				}
			} else {
				log.Printf("Received:%s\n", cmdData)
				socketToDataCommand <- cmdData
			}
		}
	}()

	go func() {
		for {
			cmdData := <-socketToDataCommand

			root := Node{}
			err := xml.Unmarshal(cmdData, &root)
			if err != nil {
				log.Printf("error parsing command: %v\n", err)
			}

			for _, event := range root.Nodes {
				dataToNode <- event
			}
		}
	}()

	go func() {
		for {
			event := <-dataToNode
			rootTag := event.XMLName.Local

			handler, exist := nc.handlers[rootTag]
			if exist {
				h := handler.Unmarshal(event)
				nodeToHanlderChannel <- h
			} else {
				log.Printf("Unknown handler: %s\n", rootTag)
			}
		}
	}()

	go func() {
		for {
			handler := <-nodeToHanlderChannel
			handler.Handle()
		}
	}()
}
