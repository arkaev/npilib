package client

import (
	"encoding/xml"
	"strings"
)

type Node struct {
	XMLName    xml.Name
	Attributes map[string]string
	Data       string
	Nodes      []*Node
}

func (e *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var nodes []*Node
	var done bool
	for !done {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch t := t.(type) {
		case xml.CharData:
			e.Data = strings.TrimSpace(string(t))
		case xml.StartElement:
			e := &Node{}
			e.UnmarshalXML(d, t)
			nodes = append(nodes, e)
		case xml.EndElement:
			done = true
		}
	}
	e.XMLName = start.Name
	e.Nodes = nodes

	e.Attributes = make(map[string]string)
	for _, attr := range start.Attr {
		e.Attributes[attr.Name.Local] = attr.Value
	}
	return nil
}

func (e *Node) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name = e.XMLName

	attrsMap := e.Attributes
	attrsArr := make([]xml.Attr, 0, len(attrsMap))
	for name, value := range attrsMap {
		attr := xml.Attr{Name: xml.Name{Space: name, Local: name}, Value: value}
		attrsArr = append(attrsArr, attr)
	}
	start.Attr = attrsArr

	return enc.EncodeElement(struct {
		Data  string `xml:",chardata"`
		Nodes []*Node
	}{
		Data:  e.Data,
		Nodes: e.Nodes,
	}, start)
}

func ParseCommand(cmd string) (Node, error) {
	n := Node{}
	err := xml.Unmarshal([]byte(cmd), &n)
	if err != nil {
		return n, err
	}

	return n, nil
}
