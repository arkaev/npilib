package npilib

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

	log.Println("Connected to socket")

	commandToSocket := make(chan NCCCommand)
	handlers := make(map[string]Handler)

	handlers["Event"] = &DoNothingHandler{}
	handlers["Command"] = &DoNothingHandler{}

	handlers["DialPlan"] = &DoNothingHandler{}
	handlers["DialplanUploadResult"] = &DoNothingHandler{}

	handlers["Success"] = &DoNothingHandler{}
	handlers["Failure"] = &DoNothingHandler{}

	handlers["FullCallsList"] = &FullCallsListHandler{}
	handlers["FullBuddyList"] = &FullBuddyListHandler{}
	handlers["ShortBuddyList"] = &DoNothingHandler{}
	handlers["BuddyListDiff"] = &DoNothingHandler{}

	handlers["LicenseUsage"] = &DoNothingHandler{}
	handlers["Progress"] = &DoNothingHandler{}

	client := RegistrationInfo{}
	responseHandlers := make(map[string]Handler)
	responseHandlers["RegisterPeer"] = &RegisterPeerHandler{config: &client, out: commandToSocket}
	responseHandlers["Register"] = &RegisterHandler{out: commandToSocket}
	responseHandlers["Subscribe"] = &DoNothingHandler{}
	handlers["Response"] = &CommonTagHandler{handlers: responseHandlers}

	requestHandlers := make(map[string]Handler)
	requestHandlers["Authenticate"] = &AuthenificateHandler{digest: auth.MD5, out: commandToSocket}
	requestHandlers["Echo"] = &EchoHandler{out: commandToSocket}
	handlers["Request"] = &CommonTagHandler{handlers: requestHandlers}

	startSender(conn, commandToSocket)
	startReceiver(conn, handlers)

	commandToSocket <- RegisterPeerCommand(auth)

	return conn, nil
}
