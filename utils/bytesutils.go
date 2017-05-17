package utils

func IntTo2byte(value int) []byte {
	ret := make([]byte, 2)
	ret[0] = byte((value >> 8) & 0xFF)
	ret[1] = byte(value & 0xFF)
	return ret
}
func IntTo4byte(value int) []byte {
	ret := make([]byte, 4)
	ret[0] = byte((value >> 24) & 0xFF)
	ret[1] = byte(value >> 16 & 0xFF)
	ret[2] = byte(value >> 8 & 0xFF)
	ret[3] = byte(value & 0xFF)
	return ret
}
func Bytes2ToInt(buff []byte) int {
	if len(buff) != 2 {
		return 0
	}
	return int(buff[0])<<8 | int(buff[1])
}
