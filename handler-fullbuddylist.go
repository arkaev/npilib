package npilib

import (
	"encoding/xml"

	c "github.com/arkaev/npilib/commands"
)

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
