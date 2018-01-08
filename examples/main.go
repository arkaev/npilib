package main

import (
	"log"
	"time"

	"github.com/arkaev/npilib"
	c "github.com/arkaev/npilib/commands"
)

const (
	url     = "docker72:3242"
	name    = "naucrm"
	keyFile = "../config/key.service." + name + ".xml"
)

func main() {
	// f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalln("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Started")

	opts := npilib.Options{
		RegisteredCB: func(nc *npilib.Conn) {
			log.Println("Successful registration")
		},
		ClosedCB: func(nc *npilib.Conn) {
			log.Println("Closed connection")
		},
	}

	nc, err := npilib.Connect(url, opts)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	nc.Subscribe("Response:Register", func(msg *npilib.Msg) {
		nc.Publish(c.CreateSubscribeRq("callslist"))
		nc.Publish(c.CreateSubscribeRq("buddylist"))
	})

	nc.Register(keyFile)

	time.Sleep(5 * time.Second)

	log.Println("Exit")
}
