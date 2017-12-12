package client

//NCCCommand base interface for response to socket
type NCCCommand interface {
}

func RegisterPeerCommand(auth *Auth) NCCCommand {
	type Params struct {
		Login       string `xml:"login,attr"`
		MaxProtocol int    `xml:"max_protocol,attr"`
		MinProtocol int    `xml:"min_protocol,attr"`
		Role        string `xml:"role,attr"`
	}

	type Request struct {
		Name   string `xml:"name,attr"`
		Params []Params
	}

	type NCCN struct {
		Request Request
	}

	return NCCN{
		Request: Request{Name: "RegisterPeer",
			Params: []Params{
				Params{Login: auth.Login, MaxProtocol: 0, MinProtocol: 0, Role: "service"}}}}
}
