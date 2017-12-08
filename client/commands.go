package client

type NCCN struct {
	Request Request
}

type Request struct {
	Name   string `xml:"name,attr"`
	Params []Params
}

type Params struct {
	Login       string `xml:"login,attr"`
	MaxProtocol int    `xml:"max_protocol,attr"`
	MinProtocol int    `xml:"min_protocol,attr"`
	Role        string `xml:"role,attr"`
}
