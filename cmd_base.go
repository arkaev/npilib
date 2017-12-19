package npilib

//Handler process command with name
type Parser interface {
	//Unmarshal node to command pojo
	Unmarshal(data []byte) NCCCommand
}

//Handler process command with name
type Handler interface {
	//Handle command and process
	Handle(cmd NCCCommand)
}
