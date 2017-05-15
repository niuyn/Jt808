package utils

import (
	"testing"
)

func TestEncodeCBCDFromString(t *testing.T) {
	t.Error(EncodeCBCDFromString("181"))
}
func TestGetBytes(t *testing.T) {
	s := "测试"
	t.Error(GetBytesWithGBK(s))
	t.Error([]byte(s))
}
