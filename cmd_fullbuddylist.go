package npilib

import (
	"encoding/xml"
)

type fullBuddyList struct {
	Endpoint []*endpoint
	Group    []*group
}

type endpoint struct {
	Login       string
	Number      string
	Type        string
	DisplayName string
	State       *state
	Numbers     *numbers
}

type state struct {
	Reason    string
	Value     string
	Timestamp uint64
	SubState  *subState
}

type subState struct {
	Name      string
	Value     string
	Timestamp uint64
}

type numbers struct {
	Numbers []*number
}

type number struct {
	Default bool
	Value   string
}

type addresses struct {
}

type group struct {
	Login       string
	Number      string
	Type        string
	DisplayName string
	Users       []*user
}

type user struct {
	Login string
}

//FullBuddyListHandler for "FullBuddyList" command
type FullBuddyListHandler struct {
	Handler
	buddyList *fullBuddyList
}

//Handle "FullBuddyList" command
func (h *FullBuddyListHandler) Handle(cmd NCCCommand) {}

type FullBuddyListRs struct {
	NCCCommand
	XMLName       xml.Name `xml:"NCC"`
	FullBuddyList *FullBuddyListRsMain
}

type FullBuddyListRsMain struct {
	XMLName  xml.Name `xml:"FullBuddyList"`
	Endpoint []*FullBuddyListRsEndpoint
	Group    []*FullBuddyListRsGroup
}

type FullBuddyListRsEndpoint struct {
	XMLName     xml.Name `xml:"Endpoint"`
	Extensions  string   `xml:"extensions,attr"`
	Login       string   `xml:"login,attr"`
	Number      string   `xml:"number,attr"`
	Type        string   `xml:"type,attr"`
	DisplayName string   `xml:"displayname,attr"`
	Numbers     *FullBuddyListRsNumbers
	State       *FullBuddyListRsState
	Addresses   *FullBuddyListRsAddresses
}

type FullBuddyListRsState struct {
	XMLName   xml.Name `xml:"State"`
	Reason    string   `xml:"reason,attr"`
	Value     string   `xml:"value,attr"`
	Timestamp uint64   `xml:"timestamp,attr"`
	SubState  *FullBuddyListRsSubState
}

type FullBuddyListRsSubState struct {
	XMLName   xml.Name `xml:"SubState"`
	Timestamp uint64   `xml:"timestamp,attr"`
	Value     bool     `xml:"value,attr"`
	Name      string   `xml:"name,attr"`
}

type FullBuddyListRsNumbers struct {
	XMLName xml.Name `xml:"Numbers"`
	Number  []*FullBuddyListRsNumber
}

type FullBuddyListRsNumber struct {
	XMLName xml.Name `xml:"Number"`
	Default bool     `xml:"default,attr"`
	Value   string   `xml:"value,attr"`
}

type FullBuddyListRsAddresses struct {
}

type FullBuddyListRsGroup struct {
	Login       string `xml:"login,attr"`
	Number      string `xml:"number,attr"`
	DisplayName string `xml:"displayname,attr"`
	User        []*FullBuddyListRsUser
}

type FullBuddyListRsUser struct {
	Login string `xml:"login,attr"`
}

//FullBuddyListHandler for "FullBuddyList" command
type FullBuddyListParser struct {
	Parser
}

//Unmarshal "FullBuddyList" command
func (h *FullBuddyListParser) Unmarshal(data []byte) NCCCommand {
	var rs FullBuddyListRs
	xml.Unmarshal(data, &rs)

	return &rs
}
