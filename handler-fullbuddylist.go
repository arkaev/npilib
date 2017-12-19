package npilib

import (
	"encoding/xml"

	c "github.com/arkaev/npilib/commands"
)

//FullBuddyListHandler for "FullBuddyList" command
type FullBuddyListHandler struct {
	Handler
}

//Handle "FullBuddyList" command
func (h *FullBuddyListHandler) Handle(cmd c.NCCCommand) {}

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
