package npilib

import "log"

//Handler process command with name
type Handler interface {
	//Unmarshal node to command pojo
	Unmarshal(node *Node) Handler
	//Handle command and process
	Handle()
}

//CommonTagHandler if handlers have same root tag but different name-attribute
type CommonTagHandler struct {
	Handler
	handlers map[string]Handler
}

//Unmarshal command for wrapper
func (h *CommonTagHandler) Unmarshal(event *Node) Handler {
	name := event.Attributes["name"]
	handler, exist := h.handlers[name]
	if exist {
		wrapped := handler.Unmarshal(event)
		return wrapped
	}

	log.Printf("Unknown handler: %s\n", name)
	return nil
}

//Handle command for wrapper
func (h *CommonTagHandler) Handle() {
}
