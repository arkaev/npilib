package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	err := connect("docker72:3242")
	if err != nil {
		fmt.Println(err)
	}
}

func connect(address string) error {
	RegisterPeer := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<NCCN>
<Request name="RegisterPeer">
<Params login="naucrm" max_protocol="0" min_protocol="0" role="service"/>
</Request>
</NCCN>`

	fmt.Println("Started")

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	conn.Write([]byte(RegisterPeer))
	conn.Write([]byte("\000"))
	fmt.Println("Sent: " + RegisterPeer)

	status, err := bufio.NewReader(conn).ReadString(0)
	if err != nil {
		return err
	}

	fmt.Println("Received: " + status)

	return nil
}
