package npilib

import (
	"log"
	"strconv"
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

//Unmarshal "FullBuddyList" command
func (h *FullBuddyListHandler) Unmarshal(node *Node) Handler {
	endpoints := make([]*endpoint, 0)
	for _, endpointNode := range node.NodesByName("Endpoint") {
		numbersArr := make([]*number, 0)
		for _, numbersNode := range endpointNode.NodesByName("Numbers") {
			for _, numberNode := range numbersNode.NodesByName("Number") {
				num := number{
					Default: "true" == numberNode.Attributes["default"],
					Value:   numberNode.Attributes["value"],
				}
				numbersArr = append(numbersArr, &num)
			}
		}

		stateNode := endpointNode.NodesByName("State")[0]
		stateTimeStr := stateNode.Attributes["timestamp"]
		stateTime, err := strconv.ParseUint(stateTimeStr, 10, 64)
		if err != nil {
			log.Printf("Error parsing '%s'. %v\n", stateTimeStr, err)
		}

		var sState *subState

		subStateNodes := stateNode.NodesByName("SubState")
		if len(subStateNodes) > 0 {
			subStateNode := stateNode.NodesByName("SubState")[0]
			subStateTimeStr := subStateNode.Attributes["timestamp"]
			subStateTime, err := strconv.ParseUint(stateTimeStr, 10, 64)
			if err != nil {
				log.Printf("Error parsing '%s'. %v\n", subStateTimeStr, err)
			}

			sState = &subState{
				Timestamp: subStateTime,
				Value:     subStateNode.Attributes["value"],
				Name:      subStateNode.Attributes["name"],
			}
		} else {
			sState = nil
		}

		st := state{
			Reason:    stateNode.Attributes["reason"],
			Value:     stateNode.Attributes["value"],
			Timestamp: stateTime,
			SubState:  sState,
		}

		ep := endpoint{
			State:       &st,
			Login:       endpointNode.Attributes["login"],
			Number:      endpointNode.Attributes["number"],
			Type:        endpointNode.Attributes["type"],
			DisplayName: endpointNode.Attributes["displayname"],
			Numbers:     &numbers{Numbers: numbersArr},
		}
		endpoints = append(endpoints, &ep)
	}

	groups := make([]*group, 0)
	for _, grNode := range node.NodesByName("Group") {
		users := make([]*user, 0)

		for _, userNode := range grNode.NodesByName("User") {
			us := user{Login: userNode.Attributes["login"]}
			users = append(users, &us)
		}

		gr := group{Login: grNode.Attributes["login"],
			Number:      grNode.Attributes["number"],
			DisplayName: grNode.Attributes["displayname"],
			Users:       users}
		groups = append(groups, &gr)
	}

	h.buddyList = &fullBuddyList{Endpoint: endpoints, Group: groups}

	return h
}

//Handle "FullBuddyList" command
func (h *FullBuddyListHandler) Handle() {}

//FullCallsListHandler for "FullCallsList" command
type FullCallsListHandler struct {
	Handler
	timeT uint64
}

//Unmarshal "FullCallsList" command
func (h *FullCallsListHandler) Unmarshal(node *Node) Handler {
	timeStr := node.Attributes["time_t"]
	timeT, err := strconv.ParseUint(timeStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing '%s'. %v\n", timeStr, err)
	}

	h.timeT = timeT
	return h
}

//Handle "FullCallsList" command
func (h *FullCallsListHandler) Handle() {}
