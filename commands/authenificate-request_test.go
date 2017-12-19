package commands

import (
	"encoding/xml"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	cmd := []byte(`<?xml version="1.0" encoding="utf-8"?>
<NCCN>
  <Request name="Authenticate">
    <Params algorithm="md5" auth_scheme="digest" method="REGISTER" nonce="CFE584AC0504F2D3000A9FBF70AB1D80" realm="ncc.net" uri="ncc.net" username="naucrm"/>
  </Request>
</NCCN>`)

	var rq AuthentificateRq
	xml.Unmarshal(cmd, &rq)

	parsed := rq.Request.Params

	if "md5" != parsed.Algorithm ||
		"digest" != parsed.AuthScheme ||
		"REGISTER" != parsed.Method ||
		"CFE584AC0504F2D3000A9FBF70AB1D80" != parsed.Nonce ||
		"ncc.net" != parsed.Realm ||
		"ncc.net" != parsed.URI ||
		"naucrm" != parsed.Username {
		t.Error("Broken parsing")
	}
}
