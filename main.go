package main

import (
	"fmt"
	"time"

	"github.com/arkaev/npilib/client"
)

const address string = "docker72:3242"
const name string = "naucrm"
const keyFile string = "config/key.service." + name + ".xml"

func main() {
	fmt.Println("Started")

	conn, err := client.Connect(address, keyFile)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	time.Sleep(time.Millisecond * 1000)

	fmt.Println("Exit")
}
