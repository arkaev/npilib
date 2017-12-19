package npilib

import (
	c "github.com/arkaev/npilib/commands"
)

//RegisterPeerCommand constuct
func RegisterPeerCommand(auth *auth) c.NCCCommand {
	return &c.RegisterPeerRq{
		Request: &c.RegisterPeerRqRequest{Name: "RegisterPeer",
			Params: &c.RegisterPeerRqParams{
				Login:       auth.Login,
				MaxProtocol: 0,
				MinProtocol: 0,
				Role:        "service"}}}
}

//SubscribeCommand constuct
func SubscribeCommand(list string) c.NCCCommand {
	return &c.SubscribeRq{
		Request: &c.SubscribeRqRequest{Name: "Subscribe",
			Params: []*c.SubscribeRqParams{
				&c.SubscribeRqParams{List: list, Enabled: true, Instant: true}}}}
}
