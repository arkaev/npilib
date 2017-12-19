package npilib

//DoNothingHandler bulk handler
type DoNothingHandler struct {
	Handler
}

//Unmarshal bulk
func (h *DoNothingHandler) Unmarshal(node *Node) Handler {
	return h
}

//Handle bulk
func (h *DoNothingHandler) Handle() {
}
