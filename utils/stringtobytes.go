package utils

func EncodeCBCDByte(str string) byte {
ret := byte(0)
if len(str) != 2 {
//return 0
} else {
for _, v := range []byte(str) {
ret = ret << 4
if v >= '0' && v <= '9' {
ret = ret | (v - '0')
} else if v >= 'A' && v <= 'F' {
ret = ret | (v - 'A' + 10)
} else if v >= 'a' && v <= 'f' {
ret = ret | (v - 'a' + 10)
}
}
}

return ret
}

func EncodeCBCDFromString(str string) []byte {
empty := make([]byte, 0)

if len(str)%2 != 0 {
return empty
}

ret := make([]byte, len(str)/2)
for i := 0; i < len(str)/2; i++ {
index := 2 * i
ret[i] = EncodeCBCDByte(str[index : index+2])

}

return ret

}