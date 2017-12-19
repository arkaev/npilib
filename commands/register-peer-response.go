package commands

import (
	"encoding/xml"
)

type RegisterPeerRs struct {
	NCCCommand
	XMLName  xml.Name `xml:"NCCN"`
	Response *RegisterPeerRsResponse
}

type RegisterPeerRsResponse struct {
	Params *RegisterPeerRsParams
}

type RegisterPeerRsParams struct {
	AllowEncoding   string `xml:"allow_encoding,attr"`
	Domain          string `xml:"domain,attr"`
	Node            string `xml:"node,attr"`
	Peer            string `xml:"peer,attr"`
	ProtocolVersion int    `xml:"protocol_version,attr"`
}
