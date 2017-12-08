package client

import (
	"fmt"
	"net"
)

const delimeter byte = 0

//Connect: create connection by address and keyFile
func Connect(address string, keyFile string) (net.Conn, error) {

	md5, err := GetDigest(keyFile)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Digest: " + md5)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected")

	var outChannel = make(chan NCCN)

	go Sender(conn, outChannel)
	go Receiver(conn)

	outChannel <- registerPeer()

	return conn, nil
}

func registerPeer() NCCN {
	return NCCN{
		Request: Request{Name: "RegisterPeer",
			Params: []Params{
				Params{Login: "naucrm", MaxProtocol: 0, MinProtocol: 0, Role: "service"}}}}
}
