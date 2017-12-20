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
	msgToHanlderChannel := make(chan *Msg)

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

			subs, exist := nc.subs[msg.Subject]
			if !exist || len(subs) == 0 {
				log.Printf("No handlers found. Skipped command: %s\n", msg.Subject)
				break
			}
			parser, exist := nc.parsers[msg.Subject]
			if exist {
				if parser != nil {
					msg.Parsed = parser.Unmarshal(msg.Data)
				}
				msg.Data = nil
				msgToHanlderChannel <- msg
			} else {
				log.Printf("Parser not found for command: %s\n", msg.Subject)
			}
		}
	}()

	go func() {
		for {
			msg := <-msgToHanlderChannel

			subs, exist := nc.subs[msg.Subject]
			if exist {
				for _, handler := range subs {
					handler(msg)
				}
			} else {
				log.Printf("Unknown handler: %s\n", msg.Subject)
			}
		}
	}()
}

func recognizeMessage(data []byte) *Msg {
	dataReader := bytes.NewReader(data)
	decoder := xml.NewDecoder(dataReader)

	cmd := &Msg{Data: data}
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
							cmd.Subject = se.Name.Local + ":" + attr.Value
							return cmd
						}
					}

					// if command has no name
					return nil
				}

				cmd.Subject = se.Name.Local
				return cmd
			}

			return nil
		}
	}

	return cmd
}
