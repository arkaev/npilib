package main

import (
	"log"
	"time"

	"github.com/arkaev/npilib/client"
)

const address string = "docker72:3242"
const name string = "naucrm"
const keyFile string = "config/key.service." + name + ".xml"

func main() {
	// f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalln("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Started")

	conn, err := client.Connect(address, keyFile)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	time.Sleep(time.Millisecond * 5000)

	log.Println("Exit")
}
