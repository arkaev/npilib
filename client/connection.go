package client

import (
	"log"
	"net"
)

const delimeter byte = 0

//RegistrationInfo contains bus registration info
type RegistrationInfo struct {
	AllowEncoding   string
	Domain          string
	Node            string
	Peer            string
	ProtocolVersion int
}

//Connect create connection by address and keyFile
func Connect(address string, keyFile string) (net.Conn, error) {

	auth, err := getAuthData(keyFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Digest: " + auth.MD5)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	log.Println("Connected")

	objectToSocket := make(chan NCCCommand)
	handlers := make(map[string]Handler)

	handlers["Echo"] = &EchoHandler{out: objectToSocket}
	handlers["Authenticate"] = &AuthenificateHandler{digest: auth.MD5, out: objectToSocket}
	client := RegistrationInfo{}
	handlers["RegisterPeer"] = &RegisterPeerHandler{config: &client, out: objectToSocket}
	handlers["Register"] = &RegisterHandler{out: objectToSocket}
	handlers["Subscribe"] = &DoNothingHandler{}

	go Sender(conn, objectToSocket)
	go Receiver(conn, handlers)

	objectToSocket <- RegisterPeerCommand(auth)

	return conn, nil
}
