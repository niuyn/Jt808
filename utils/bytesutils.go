package utils

func IntTo2byte(value int) []byte {
	ret := make([]byte, 2)
	ret[0] = byte((value >> 8) & 0xFF)
	ret[1] = byte(value & 0xFF)
	return ret
}
