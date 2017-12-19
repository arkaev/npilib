package npilib

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"io"
	"log"
	"time"

	c "github.com/arkaev/npilib/commands"
)

//Receiver for commands from socket
func startReceiver(nc *Conn) {
	socketToDataCommand := make(chan []byte)
	msgToHanlderChannel := make(chan *CommandWrapper)

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
			msg := recognizeMessage(cmdData)

			_, exist := nc.handlers[msg.Command]
			if !exist {
				log.Printf("No handler found. Skipped command: %s\n", msg.Command)
				break
			}
			parser, exist := nc.parsers[msg.Command]
			if exist {
				if parser != nil {
					msg.Parsed = parser.Unmarshal(msg.Data)
				}
				msg.Data = nil
				msgToHanlderChannel <- msg
			} else {
				log.Printf("Parser not found for command: %s\n", msg.Command)
			}
		}
	}()

	go func() {
		for {
			msg := <-msgToHanlderChannel

			handler, exist := nc.handlers[msg.Command]
			if exist {
				handler.Handle(msg.Parsed)
			} else {
				log.Printf("Unknown handler: %s\n", msg.Command)
			}
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
	Parsed  c.NCCCommand
}

func recognizeMessage(data []byte) *CommandWrapper {
	dataReader := bytes.NewReader(data)
	decoder := xml.NewDecoder(dataReader)

	cmd := &CommandWrapper{Data: data}
	level := 0
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			level++
			if level == 1 {
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
				} else {
					return nil
				}
			}

			if level == 2 {
				if se.Name.Local == "Request" || se.Name.Local == "Response" {
					for _, attr := range se.Attr {
						if attr.Name.Local == "name" {
							cmd.Command = se.Name.Local + ":" + attr.Value
							return cmd
						}
					}

					// if command has no name
					return nil
				}

				cmd.Command = se.Name.Local
				return cmd
			}

			return nil
		}
	}

	return cmd
}
