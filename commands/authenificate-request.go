package commands

import "encoding/xml"

// AuthenificateRq structure contains "Authenificate" request data
// type AuthenificateRq struct {
// 	NCCCommand
// 	Algorithm  string
// 	AuthScheme string
// 	Method     string
// 	Nonce      string
// 	Realm      string
// 	URI        string
// 	Username   string
// }

type AuthentificateRq struct {
	XMLName xml.Name `xml:"NCCN"`
	Request *AuthentificateRqRequest
}

type AuthentificateRqRequest struct {
	Name   string `xml:"name,attr"`
	Params *AuthentificateRqParams
}

type AuthentificateRqParams struct {
	Algorithm  string `xml:"algorithm,attr"`
	AuthScheme string `xml:"auth_scheme,attr"`
	Method     string `xml:"method,attr"`
	Nonce      string `xml:"nonce,attr"`
	Realm      string `xml:"realm,attr"`
	URI        string `xml:"uri,attr"`
	Username   string `xml:"username,attr"`
}
