package npilib

import (
	"testing"
)

func TestCalculateMD5(t *testing.T) {
	val := calculateMD5("c4401cb66d13b3fd6f5ab25be0503123", "86165664D8F9020AD2CAC861928E2AA7", "REGISTER", "ncc.net")
	if "d68c9f8387a3249f0e9a6e02d080abe4" != val {
		t.Error("Wrong result")
	}
}
