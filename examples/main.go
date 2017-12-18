package main

import (
	"log"
	"time"

	"github.com/arkaev/npilib"
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

	opts := npilib.Options{KeyFile: keyFile}

	conn, err := npilib.Connect(url, opts)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	time.Sleep(time.Millisecond * 1000)

	log.Println("Exit")
}
