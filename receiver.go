package npilib

import (
	"bufio"
	"bytes"
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
				log.Printf("Received:\n%s\n", cmdData)
				socketToDataCommand <- cmdData
			}
		}
	}()

	go func() {
		for {
			cmdData := <-socketToDataCommand
			cmd := recognizeMessage(cmdData)

			log.Printf("Recognized command: %s\n", cmd.Command)

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

// CommandWrapper contains base information about command
type CommandWrapper struct {
	NCC     string
	Command string
	From    string
	To      string
	Data    []byte
}

func recognizeMessage(data []byte) *CommandWrapper {
	dataReader := bytes.NewReader(data)
	decoder := xml.NewDecoder(dataReader)

	cmd := &CommandWrapper{Data: data}

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "NCC" || se.Name.Local == "NCCN" {
				cmd.NCC = se.Name.Local
				for _, attr := range se.Attr {
					if attr.Name.Local == "from" {
						cmd.From = attr.Value
					} else if attr.Name.Local == "to" {
						cmd.To = attr.Value
					}
				}
				break
			} else if se.Name.Local == "Request" || se.Name.Local == "Response" {
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						cmd.Command = attr.Value
						return cmd
					}
				}

				// if command has no name
				return nil
			} else {
				cmd.Command = se.Name.Local
				return cmd
			}
		}
	}

	return cmd
}
