package client

import (
	"bufio"
	"encoding/xml"
	"io"
	"log"
	"net"
	"time"
)

//Receiver for commands from socket
func Receiver(conn net.Conn, handlers map[string]Handler) {
	socketToParser := make(chan string)

	go commandParser(socketToParser, handlers)

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

func commandParser(socketToParser <-chan string, handlers map[string]Handler) {
	for {
		cmdStr := <-socketToParser
		root, err := parseCommand(cmdStr)
		if err != nil {
			log.Printf("error parsing command: %v\n", err)
			return
		}

		for _, request := range root.Nodes {
			commandName := request.Attributes["name"]

			handler, exist := handlers[commandName]
			if exist {
				handler.Unmarshal(request)
				handler.Handle()
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
