package client

import (
	"log"
	"net"
)

const delimeter byte = 0

//Connect create connection by address and keyFile
func Connect(address string, keyFile string) (net.Conn, error) {

	auth, err := GetAuthData(keyFile)
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
	handlers["Echo"] = &EchoHandler{outCommands: objectToSocket}
	handlers["Authenticate"] = &AuthenificateHandler{digest: auth.MD5, outCommands: objectToSocket}
	handlers["RegisterPeer"] = &RegisterPeerHandler{outCommands: objectToSocket}
	handlers["Register"] = &RegisterHandler{}

	go Sender(conn, objectToSocket)

	receiver := &Receiver{handlers}
	go receiver.Start(conn)

	objectToSocket <- registerPeer(auth)

	return conn, nil
}

func registerPeer(auth *Auth) NCCCommand {
	type Params struct {
		Login       string `xml:"login,attr"`
		MaxProtocol int    `xml:"max_protocol,attr"`
		MinProtocol int    `xml:"min_protocol,attr"`
		Role        string `xml:"role,attr"`
	}

	type Request struct {
		Name   string `xml:"name,attr"`
		Params []Params
	}

	type NCCN struct {
		Request Request
	}

	return NCCN{
		Request: Request{Name: "RegisterPeer",
			Params: []Params{
				Params{Login: auth.Login, MaxProtocol: 0, MinProtocol: 0, Role: "service"}}}}
}
