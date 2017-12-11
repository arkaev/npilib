package client

import (
	"encoding/xml"
)

func ParseCommand(cmd string) (Node, error) {
	n := Node{}
	err := xml.Unmarshal([]byte(cmd), &n)
	if err != nil {
		return n, err
	}

	return n, nil
}
