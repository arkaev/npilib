package client

import (
	"encoding/xml"
	"fmt"
	"net"
)

func Sender(conn net.Conn, commands <-chan *Node) {
	for {
		cmd := <-commands
		output, err := xml.MarshalIndent(cmd, "", "    ")
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}

		conn.Write(output)
		conn.Write([]byte{delimeter})
		fmt.Print("Sent: ")
		fmt.Println(string(output))
	}
}
