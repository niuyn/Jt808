package model

import (
	"encoding/hex"
	"testing"
)

func TestTerminalGpsBody_GetBytes(t *testing.T) {
	dev := &TerminalGpsBody{}
	dev.SetTime("170515162830")
	t.Error(hex.EncodeToString(dev.GetBytes()))
}
