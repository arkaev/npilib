package npilib

import c "github.com/arkaev/npilib/commands"

//Handler process command with name
type Parser interface {
	//Unmarshal node to command pojo
	Unmarshal(data []byte) c.NCCCommand
}
