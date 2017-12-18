package npilib

import (
	"testing"
)

func TestCalculateMD5(t *testing.T) {
	val := calculateMD5(digestTemp, "FEEABF56889D7817EF3DEA4B5F631CAC", "REGISTER", "ncc.net")
	if "a7c8c5dc9ac454e8e2c4ff0102ea3fbd" != val {
		t.Error("Wrong result")
	}
}
