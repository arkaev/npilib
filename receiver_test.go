package npilib

import "testing"

func TestRecognizeNCCNRequest(t *testing.T) {
	cmd := []byte(`<?xml version="1.0" encoding="utf-8"?>
<NCCN>
  <Request name="Authenticate">
    <Params algorithm="md5" auth_scheme="digest" method="REGISTER" nonce="CFE584AC0504F2D3000A9FBF70AB1D80" realm="ncc.net" uri="ncc.net" username="naucrm"/>
  </Request>
</NCCN>`)

	msg := recognizeMessage(cmd)
	if msg.Subject != "Request:Authenticate" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}

func TestRecognizeNCCNResponse(t *testing.T) {
	cmd := []byte(`<?xml version="1.0" encoding="utf-8"?>
<NCCN>
  <Response name="RegisterPeer">
    <Params allow_encoding="ncc3, ncca2, bzip2, lzma, lz4" domain="domain" node="node" peer="naucrm-191" protocol_version="100"/>
  </Response>
</NCCN>`)

	msg := recognizeMessage(cmd)
	if msg.Subject != "Response:RegisterPeer" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}

func TestRecognizeNCCResponse(t *testing.T) {
	cmd := []byte(`<NCC from="naubuddy-20.node.domain" to="naucrm-191.node.domain">
<Response name="Subscribe">
  <Params/>
 </Response></NCC>`)

	msg := recognizeMessage(cmd)
	if msg.Subject != "Response:Subscribe" ||
		msg.From != "naubuddy-20.node.domain" ||
		msg.To != "naucrm-191.node.domain" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}

func TestRecognizeCommandWithTag(t *testing.T) {
	cmd := []byte(`<NCC>
	<FullBuddyList name="NotUsed">
		<SomeInner/>
	</FullBuddyList>
</NCC>`)

	msg := recognizeMessage(cmd)
	if msg.Subject != "FullBuddyList" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}

func TestRecognizeWrongCommand(t *testing.T) {
	cmd := []byte(`<Wrong>
	<FullBuddyList name="NotUsed">
		<SomeInner/>
	</FullBuddyList>
 </Wrong>`)

	msg := recognizeMessage(cmd)
	if msg != nil {
		t.Error("Wrong recognized")
	}
}

func TestRecognizeWrongRequestWithoutName(t *testing.T) {
	cmd := []byte(`<NCC>
	<Request>
		<SomeInner/>
	</Request>
 </NCC>`)

	msg := recognizeMessage(cmd)
	if msg != nil {
		t.Error("Wrong recognized")
	}
}
