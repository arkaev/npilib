package commands

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestFullBaddyListParse(t *testing.T) {
	cmd := []byte(`<NCC from="naubuddy-20.node.domain" to="naucrm-194.node.domain">
		<FullBuddyList>
			<Endpoint login="agent" number="agent" type="phone" displayname="Иванов Иван Иванович">
				<Numbers>
					<Number default="true" value="agent"/>
				</Numbers>
				<State reason="some_reason" value="offline" timestamp="0"/>
				<Addresses/>
			</Endpoint>
			<Endpoint login="ivr" number="0001" type="ivr" displayname="NauIVR">
				<Numbers>
					<Number value="0001"/>
					<Number default="true" value="0002"/>
					<Number value="ivr"/>
				</Numbers>
				<State reason="some_reason" value="available" timestamp="1512709577">
					<SubState timestamp="1512709577" value="true" name="normal"/>
					<SubState timestamp="123" value="false" name="wrapup"/>
				</State>
			</Endpoint>
			<Group login="root" number="root" displayname="Тестовая компания">
				<User login="supervisor"/>
				<User login="servers"/>
				<User login="agent"/>
			</Group>
			<Group login="servers" number="servers" displayname="Services">
				<User login="nauss_0"/>
				<User login="ivr"/>
			</Group>
		</FullBuddyList>
	</NCC>`)

	var rs FullBuddyListRs
	xml.Unmarshal(cmd, &rs)

	assertEqual(t, rs.FullBuddyList.Endpoint[1].Login, "ivr")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Number, "0001")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Type, "ivr")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].DisplayName, "NauIVR")

	assertEqual(t, rs.FullBuddyList.Endpoint[1].Numbers.Number[0].Default, false)
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Numbers.Number[0].Value, "0001")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Numbers.Number[1].Default, true)
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Numbers.Number[1].Value, "0002")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Numbers.Number[2].Default, false)
	assertEqual(t, rs.FullBuddyList.Endpoint[1].Numbers.Number[2].Value, "ivr")

	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.Reason, "some_reason")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.Value, "available")
	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.Timestamp, uint64(1512709577))

	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.SubState[0].Timestamp, uint64(1512709577))
	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.SubState[0].Value, true)
	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.SubState[0].Name, "normal")

	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.SubState[1].Timestamp, uint64(123))
	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.SubState[1].Value, false)
	assertEqual(t, rs.FullBuddyList.Endpoint[1].State.SubState[1].Name, "wrapup")

	assertEqual(t, rs.FullBuddyList.Group[1].Login, "servers")
	assertEqual(t, rs.FullBuddyList.Group[1].Number, "servers")
	assertEqual(t, rs.FullBuddyList.Group[1].DisplayName, "Services")

	assertEqual(t, rs.FullBuddyList.Group[1].User[0].Login, "nauss_0")
	assertEqual(t, rs.FullBuddyList.Group[1].User[1].Login, "ivr")
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal(fmt.Sprintf("Assert failed: %v != %v", a, b))
}
