package commands

import (
	"encoding/xml"
)

type EchoRq struct {
	NCCCommand
	XMLName xml.Name `xml:"NCC"`
	Request *EchoRqRequest
}

type EchoRqRequest struct {
	Name string `xml:"name,attr"`
}
