package client

import (
	"bufio"
	"fmt"
	"net"
)

func Receiver(conn net.Conn) {
	var inQueue = make(chan string)
	go commandParser(inQueue)

	bufReader := bufio.NewReader(conn)
	for {
		msg, err := bufReader.ReadString(delimeter)
		if err != nil {
			break
		}
		fmt.Print("Received: ")
		fmt.Println(msg)
		inQueue <- msg
	}
}

func commandParser(strCommands <-chan string) {
	cmdStr := <-strCommands
	root, err := ParseCommand(cmdStr)
	if err != nil {
		fmt.Printf("error parsing command: %v\n", err)
		return
	}

	for _, Request := range root.Nodes {
		commandName := Request.Attributes["name"]

		switch commandName {
		// case "Authenticate":
		//HandleAuthenificate()
		default:
			fmt.Printf("Unknown command: %s\n", commandName)
		}
	}
}
