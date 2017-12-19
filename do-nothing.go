package npilib

import c "github.com/arkaev/npilib/commands"

//DoNothingHandler bulk handler
type DoNothingHandler struct {
	Handler
}

//Handle bulk
func (h *DoNothingHandler) Handle(cmd c.NCCCommand) {
}
