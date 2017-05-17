package utils

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

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
	if len(str) == 0 {
		return empty
	}
	// 奇数 前面补零
	if len(str)%2 != 0 {
		str = "0" + str
	}

	ret := make([]byte, len(str)/2)
	for i := 0; i < len(str)/2; i++ {
		index := 2 * i
		ret[i] = EncodeCBCDByte(str[index : index+2])

	}

	return ret

}

// 获得string 的GBK　编码
func GetBytesWithGBK(str string) []byte {
	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(str)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return nil
	}
	return data
}

func GetStringWithGBK(buffer []byte) string {
	strBytes, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(buffer), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		return ""
	}
	return string(strBytes)

}
