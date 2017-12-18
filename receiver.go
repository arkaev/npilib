package npilib

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"time"
)

//Receiver for commands from socket
func startReceiver(conn *Conn, handlers map[string]Handler) {
	socketToStrCommand := make(chan string)
	strToNode := make(chan *Node)
	nodeToHanlderChannel := make(chan Handler)

	go func() {
		bufReader := bufio.NewReader(conn.conn)
		for {
			msg, err := bufReader.ReadString(delimeter)
			if err != nil {
				if err == io.EOF {
					//sleep if no data
					time.Sleep(time.Millisecond * 10)
				} else {
					log.Println("Unexpected read error: ", err)
					break
				}
			} else {
				log.Println("Received:\n" + msg)
				socketToStrCommand <- msg
			}
		}
	}()

	go func() {
		for {
			cmdStr := <-socketToStrCommand

			root := Node{}
			err := xml.Unmarshal([]byte(cmdStr), &root)
			if err != nil {
				log.Printf("error parsing command: %v\n", err)
			}

			for _, event := range root.Nodes {
				strToNode <- event
			}
		}
	}()

	go func() {
		for {
			event := <-strToNode
			rootTag := event.XMLName.Local

			handler, exist := handlers[rootTag]
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
