package npilib

import "testing"

const (
	cmdAuthenticate = `<?xml version="1.0" encoding="utf-8"?>
<NCCN>
  <Request name="Authenticate">
    <Params algorithm="md5" auth_scheme="digest" method="REGISTER" nonce="CFE584AC0504F2D3000A9FBF70AB1D80" realm="ncc.net" uri="ncc.net" username="naucrm"/>
  </Request>
</NCCN>`

	cmdRegisterPeer = `<?xml version="1.0" encoding="utf-8"?>
<NCCN>
  <Response name="RegisterPeer">
    <Params allow_encoding="ncc3, ncca2, bzip2, lzma, lz4" domain="domain" node="node" peer="naucrm-191" protocol_version="100"/>
  </Response>
</NCCN>`

	cmdSubscribe = `<NCC from="naubuddy-20.node.domain" to="naucrm-191.node.domain">
<Response name="Subscribe">
  <Params/>
 </Response></NCC>`
)

func TestRecognizeNCCNRequest(t *testing.T) {
	msg := recognizeMessage([]byte(cmdAuthenticate))
	if msg.NCC != "NCCN" ||
		msg.Command != "Authenticate" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}

func TestRecognizeNCCNResponse(t *testing.T) {
	msg := recognizeMessage([]byte(cmdRegisterPeer))
	if msg.NCC != "NCCN" ||
		msg.Command != "RegisterPeer" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}

func TestRecognizeNCCResponse(t *testing.T) {
	msg := recognizeMessage([]byte(cmdSubscribe))
	if msg.NCC != "NCC" ||
		msg.Command != "Subscribe" ||
		msg.From != "naubuddy-20.node.domain" ||
		msg.To != "naucrm-191.node.domain" ||
		msg.Data == nil {
		t.Error("Not recognized")
	}
}
