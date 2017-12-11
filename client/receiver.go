package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func Receiver(conn net.Conn, outCommands chan<- *Node) {
	var inQueue = make(chan string)
	go commandParser(inQueue, outCommands)

	bufReader := bufio.NewReader(conn)
	for {
		msg, err := bufReader.ReadString(delimeter)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Unexpected read error: ", err)
				break
			}
		} else {
			fmt.Print("Received: ")
			fmt.Println(msg)
			inQueue <- msg
		}
	}
}

func commandParser(strCommands <-chan string, outCommands chan<- *Node) {
	cmdStr := <-strCommands
	root, err := ParseCommand(cmdStr)
	if err != nil {
		fmt.Printf("error parsing command: %v\n", err)
		return
	}

	for _, request := range root.Nodes {
		commandName := request.Attributes["name"]

		var err error

		switch commandName {
		case "Authenticate":
			HandleAuthenificate(request, outCommands)
		default:
			fmt.Printf("Unknown command: %s\n", commandName)
		}

		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}
}
