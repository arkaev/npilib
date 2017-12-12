package client

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"net"
	"time"
)

type Receiver struct {
	handlers map[string]Handler
}

//Start receiving commands from socket
func (r *Receiver) Start(conn net.Conn) {
	socketToParser := make(chan string)

	parser := &commandParser{handlers: r.handlers}
	go parser.parse(socketToParser)

	bufReader := bufio.NewReader(conn)
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
			socketToParser <- msg
		}
	}
}

type commandParser struct {
	handlers map[string]Handler
}

func (p *commandParser) parse(socketToParser <-chan string) {
	for {
		cmdStr := <-socketToParser
		root, err := parseCommand(cmdStr)
		if err != nil {
			log.Printf("error parsing command: %v\n", err)
			return
		}

		for _, request := range root.Nodes {
			commandName := request.Attributes["name"]

			handler, exist := p.handlers[commandName]
			if exist {
				handler.Process(request)
			} else {
				log.Printf("Unknown command: %s\n", commandName)
			}
		}
	}
}

func parseCommand(cmd string) (Node, error) {
	n := Node{}
	err := xml.Unmarshal([]byte(cmd), &n)
	if err != nil {
		return n, err
	}

	return n, nil
}
