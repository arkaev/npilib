package npilib

import (
	"log"
	"net"
	"net/url"
)

const (
	// delimeter for messages
	delimeter byte = 0
)

// A Conn represents a bare connection to a server.
type Conn struct {
	url             *url.URL
	conn            net.Conn
	digest          string
	allowEncoding   string
	domain          string
	node            string
	peer            string
	protocolVersion int
}

// Options can be used to create a customized connection.
type Options struct {
	//KeyFile is url to key-file with authentification digest
	KeyFile string
}

//Connect create connection by address and keyFile
func Connect(url string, options Options) (*Conn, error) {

	conn := &Conn{}

	auth, err := getAuthData(options.KeyFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	conn.digest = auth.MD5
	log.Printf("Digest: %s\n", conn.digest)

	conn.conn, err = net.Dial("tcp", url)
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

	responseHandlers := make(map[string]Handler)
	responseHandlers["RegisterPeer"] = &RegisterPeerHandler{config: conn, out: commandToSocket}
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

// Send will write data to socket
func (c *Conn) Send(cmd []byte) {
	c.conn.Write(cmd)
	c.conn.Write([]byte{delimeter})
}

// Close will close the connection
func (c *Conn) Close() {
	c.conn.Close()
}
