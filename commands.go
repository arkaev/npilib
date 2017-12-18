package npilib

//NCCCommand base interface for response to socket
type NCCCommand interface {
}

//RegisterPeerCommand constuct
func RegisterPeerCommand(auth *auth) NCCCommand {
	type Params struct {
		Login       string `xml:"login,attr"`
		MaxProtocol int    `xml:"max_protocol,attr"`
		MinProtocol int    `xml:"min_protocol,attr"`
		Role        string `xml:"role,attr"`
	}

	type Request struct {
		Name   string `xml:"name,attr"`
		Params []*Params
	}

	type NCCN struct {
		Request *Request
	}

	return &NCCN{
		Request: &Request{Name: "RegisterPeer",
			Params: []*Params{
				&Params{Login: auth.Login, MaxProtocol: 0, MinProtocol: 0, Role: "service"}}}}
}

//SubscribeCommand constuct
func SubscribeCommand(list string) NCCCommand {
	type Params struct {
		List    string `xml:"list,attr"`
		Enabled bool   `xml:"enabled,attr"`
		Instant bool   `xml:"instant,attr"`
	}

	type Request struct {
		Name   string `xml:"name,attr"`
		Params []*Params
	}

	type NCC struct {
		Request *Request
	}

	return &NCC{
		Request: &Request{Name: "Subscribe",
			Params: []*Params{
				&Params{List: list, Enabled: true, Instant: true}}}}
}
