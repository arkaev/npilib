package npilib

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	cmd := []byte(`<?xml version="1.0" encoding="utf-8"?>
<NCCN>
  <Request name="Authenticate">
    <Params algorithm="md5" auth_scheme="digest" method="REGISTER" nonce="CFE584AC0504F2D3000A9FBF70AB1D80" realm="ncc.net" uri="ncc.net" username="naucrm"/>
  </Request>
</NCCN>`)

	parser := &AuthenificateRqParser{}
	parsed := parser.Unmarshal(cmd).(*AuthenificateRq)

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

func TestCalculateMD5(t *testing.T) {
	val := calculateMD5("c4401cb66d13b3fd6f5ab25be0503123", "86165664D8F9020AD2CAC861928E2AA7", "REGISTER", "ncc.net")
	if "d68c9f8387a3249f0e9a6e02d080abe4" != val {
		t.Error("Wrong result")
	}
}
