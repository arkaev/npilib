package npilib

import (
	"log"
	"net"
	"net/url"
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
	subs            map[int64]*Subscription
}

// Options can be used to create a customized connection.
type Options struct {
	//KeyFile is url to key-file with authentification digest
	KeyFile string
}

// A Subscription represents interest in a given subject.
type Subscription struct {
	sid int64

	// Subject that represents this subscription
	Subject string

	// Optional queue group name. If present, all subscriptions with the
	// same name will form a distributed queue, and each message will
	// only be processed by one member of the group.
	Queue string

	conn *Conn
	mcb  MsgHandler
}

// Msg is a structure used by Subscribers and PublishMsg().
type Msg struct {
	Subject string
	Reply   string
	Data    []byte
	Sub     *Subscription
	next    *Msg
}

// MsgHandler is a callback function that processes messages delivered to
// asynchronous subscribers.
type MsgHandler func(msg *Msg)

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

	handlers["Response"] = &CommonTagHandler{
		handlers: map[string]Handler{
			"RegisterPeer": &RegisterPeerHandler{config: conn, out: commandToSocket},
			"Register":     &RegisterHandler{out: commandToSocket},
			"Subscribe":    &DoNothingHandler{},
		}}

	handlers["Request"] = &CommonTagHandler{
		handlers: map[string]Handler{
			"Authenticate": &AuthenificateHandler{digest: auth.MD5, out: commandToSocket},
			"Echo":         &EchoHandler{out: commandToSocket},
		}}

	startSender(conn, commandToSocket)
	startReceiver(conn, handlers)

	commandToSocket <- RegisterPeerCommand(auth)

	return conn, nil
}

// subscribe is the internal subscribe function that indicates interest in a subject.
func (nc *Conn) subscribe(subj, queue string, cb MsgHandler, ch chan *Msg) (*Subscription, error) {
	sub := &Subscription{Subject: subj, Queue: queue, mcb: cb, conn: nc}

	sub.sid = nc.ssid
	nc.subs[sub.sid] = sub

	return sub, nil
}

// Subscribe will execute handler on subject event
func (nc *Conn) Subscribe(subj string, cb MsgHandler) (*Subscription, error) {
	return nc.subscribe(subj, _EMPTY_, cb, nil)
}

// Send will write data to socket
func (nc *Conn) Send(cmd []byte) {
	nc.conn.Write(cmd)
	nc.conn.Write([]byte{delimeter})
}

// Close will close the connection
func (nc *Conn) Close() {
	nc.conn.Close()
}
