package commands

// //FullCallsListHandler for "FullCallsList" command
// type FullCallsListHandler struct {
// 	Handler
// 	timeT uint64
// }

// //Unmarshal "FullCallsList" command
// func (h *FullCallsListHandler) Unmarshal(node *Node) Handler {
// 	timeStr := node.Attributes["time_t"]
// 	timeT, err := strconv.ParseUint(timeStr, 10, 64)
// 	if err != nil {
// 		log.Printf("Error parsing '%s'. %v\n", timeStr, err)
// 	}

// 	h.timeT = timeT
// 	return h
// }

// //Handle "FullCallsList" command
// func (h *FullCallsListHandler) Handle() {}
