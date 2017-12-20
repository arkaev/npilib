package npilib

import (
	"log"
	"net"
	"net/url"

	c "github.com/arkaev/npilib/commands"
)

const (
	// delimeter for messages
	delimeter byte = 0

	_EMPTY_ = ""
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
	ssid            int64
	subs            map[string]([]MsgHandler)

	options         Options
	parsers         map[string]Parser
	commandToSocket chan c.NCCCommand
}

// ConnHandler is used for asynchronous events such as
// disconnected and closed connections.
type ConnHandler func(*Conn)

// Options can be used to create a customized connection.
type Options struct {
	// ClosedCB sets the closed handler that is called when a client will
	// no longer be connected.
	ClosedCB ConnHandler

	// DisconnectedCB sets the disconnected handler that is called
	// whenever the connection is disconnected.
	DisconnectedCB ConnHandler

	// ReconnectedCB sets the reconnected handler called whenever
	// the connection is successfully reconnected.
	ReconnectedCB ConnHandler

	// RegisteredCB sets the callback that is invoked whenever
	// registered on bus
	RegisteredCB ConnHandler
}

// Msg is a structure used by Subscribers and PublishMsg().
type Msg struct {
	Subject string
	From    string
	To      string
	Data    []byte
	Parsed  c.NCCCommand
}

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandler func(msg *Msg)

//Connect create connection by address and keyFile
func Connect(url string, options Options) (*Conn, error) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to socket")

	nc := &Conn{
		conn:            conn,
		options:         options,
		commandToSocket: make(chan c.NCCCommand),
		subs:            make(map[string]([]MsgHandler)),
		parsers: map[string]Parser{
			"FullCallsList": nil,
			"FullBuddyList": &FullBuddyListParser{},

			"Response:RegisterPeer": &RegisterPeerRsParser{},
			"Response:Register":     nil,
			"Response:Subscribe":    nil,

			"Request:Echo":         nil,
			"Request:Authenticate": &AuthenificateRqParser{},
		},
	}

	nc.Subscribe("Event", func(*Msg) {})
	nc.Subscribe("Command", func(*Msg) {})

	nc.Subscribe("DialPlan", func(*Msg) {})
	nc.Subscribe("DialplanUploadResult", func(*Msg) {})

	nc.Subscribe("Success", func(*Msg) {})
	nc.Subscribe("Failure", func(*Msg) {})

	nc.Subscribe("FullCallsList", func(*Msg) {})
	nc.Subscribe("FullBuddyList", func(*Msg) {})
	nc.Subscribe("ShortBuddyList", func(*Msg) {})
	nc.Subscribe("BuddyListDiff", func(*Msg) {})

	nc.Subscribe("LicenseUsage", func(*Msg) {})
	nc.Subscribe("Progress", func(*Msg) {})

	nc.Subscribe("Response:RegisterPeer", func(msg *Msg) { HandleRegisterPeer(nc, msg) })
	nc.Subscribe("Response:Register", func(*Msg) {
		if nc.options.RegisteredCB != nil {
			nc.options.RegisteredCB(nc)
		}
	})
	nc.Subscribe("Response:Subscribe", func(*Msg) {})
	nc.Subscribe("Request:Authenticate", func(msg *Msg) { HandleAuthenificate(nc, msg) })
	nc.Subscribe("Request:Echo", func(*Msg) {
		nc.Publish(c.CreateEchoResponse())
	})

	startSender(nc)
	startReceiver(nc)

	return nc, nil
}

// Register will start process of negotiating.
func (nc *Conn) Register(keyFile string) {
	auth, err := getAuthData(keyFile)
	if err != nil {
		log.Println(err)
		return
	}

	nc.digest = auth.MD5
	log.Printf("Digest: %s\n", nc.digest)

	nc.commandToSocket <- c.CreateRegisterPeerCommand(auth.Login)
}

// Subscribe will execute handler on subject event
func (nc *Conn) Subscribe(subj string, cb MsgHandler) {
	subjectSubscriptions := nc.subs[subj]
	nc.subs[subj] = append(subjectSubscriptions, cb)
}

// Publish will send command to bus
func (nc *Conn) Publish(cmd c.NCCCommand) {
	nc.commandToSocket <- cmd
}

// Send will write data to socket
func (nc *Conn) send(cmd []byte) {
	nc.conn.Write(cmd)
	nc.conn.Write([]byte{delimeter})
}

// Close will close the connection
func (nc *Conn) Close() {
	if nc.options.ClosedCB != nil {
		nc.options.ClosedCB(nc)
	}
	nc.conn.Close()
}
