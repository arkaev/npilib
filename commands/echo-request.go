package commands

import (
	"encoding/xml"
)

type EchoRs struct {
	NCCCommand
	XMLName  xml.Name `xml:"NCC"`
	Response *EchoRsResponse
}

type EchoRsResponse struct {
	Name string `xml:"name,attr"`
}
