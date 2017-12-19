package commands

import "encoding/xml"

type AuthenticateRs struct {
	XMLName  xml.Name `xml:"NCCN"`
	Response *AuthenticateRsResponse
}

type AuthenticateRsResponse struct {
	Name  string `xml:"name,attr"`
	Param *AuthenticateRsParam
}

type AuthenticateRsParam struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
