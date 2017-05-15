package model

import "testing"

func TestGetBytes(t *testing.T) {
	dev := &TerminalRegisterBody{}
	dev.city = 1
	dev.province = 2
	t.Error(dev.GetBytes())
}
