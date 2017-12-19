package commands

import (
	"encoding/xml"
	"log"
	"testing"
)

func TestFullCallsListParse(t *testing.T) {
	cmd := []byte(`<NCC from="naubuddy-17.node.domain" to="naucrm-68.node.domain">
		<FullCallsList time_t="1513702855"/></NCC>`)

	var rs FullCallsList
	err := xml.Unmarshal(cmd, &rs)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	assertEqual(t, rs.FullCallsList.TimeT, uint64(1513702855))
}
