package npilib

//DoNothingHandler bulk handler
type DoNothingHandler struct {
	Handler
}

//Handle bulk
func (h *DoNothingHandler) Handle(cmd NCCCommand) {
}
