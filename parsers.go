package npilib

import (
	"encoding/xml"

	c "github.com/arkaev/npilib/commands"
)

type Parser interface {
	//Unmarshal node to command pojo
	Unmarshal(data []byte) c.NCCCommand
}

type AuthenificateRqParser struct {
	Parser
}

//Unmarshal "Authenificate" command
func (h *AuthenificateRqParser) Unmarshal(data []byte) c.NCCCommand {
	var auth c.AuthentificateRq
	xml.Unmarshal(data, &auth)

	return &auth
}

//FullBuddyListHandler for "FullBuddyList" command
type FullBuddyListParser struct {
	Parser
}

//Unmarshal "FullBuddyList" command
func (h *FullBuddyListParser) Unmarshal(data []byte) c.NCCCommand {
	var rs c.FullBuddyListRs
	xml.Unmarshal(data, &rs)

	return &rs
}

type RegisterPeerRsParser struct {
	Parser
}

func (p *RegisterPeerRsParser) Unmarshal(data []byte) c.NCCCommand {
	var rq c.RegisterPeerRs
	xml.Unmarshal(data, &rq)

	return &rq
}
