package JT808

import (
	//"sync/atomic"
	"Jt808/utils"
	"bytes"
	"errors"
)

// 消息包封装项

type PackageOption struct {
	totalNum  int //消息包总数
	packetSeq int // 消息包序号
}

func generateMsgBodyProps(msgLen, enctyptionType, reversed_14_15 int, isSubPackage bool) (error, int) {
	if msgLen >= 1024 {
		return errors.New("msglen is too long"), 0
	}

	// [ 0-9 ] 0000,0011,1111,1111(3FF)(消息体长度)
	// [10-12] 0001,1100,0000,0000(1C00)(加密类型)
	// [ 13_ ] 0010,0000,0000,0000(2000)(是否有子包)
	// [14-15] 1100,0000,0000,0000(C000)(保留位)
	subPkg := 0
	if isSubPackage {
		subPkg = 1
	}
	ret := (msgLen & 0x3FF) | ((enctyptionType << 10) & 0x1C00) | ((subPkg << 13) & 0x2000) | ((reversed_14_15 << 14) & 0xC000)
	return nil, ret & 0xFFFF
}

func generateMsgHeader(msgID, msgProps, sequence int, phone string) []byte {
	buff := bytes.Buffer{}
	buff.Write(utils.IntTo2byte(msgID))
	buff.Write(utils.IntTo2byte(msgProps))
	// 手机号 不足12位在前补0
	buff.Write(utils.EncodeCBCDFromString(leftPadding(12, phone)))
	buff.Write(utils.IntTo2byte(sequence))
	return buff.Bytes()
}
func generateMsgHeaderWithEncap(msgID, msgProps, sequence int, phone string, packageOption PackageOption) []byte {
	buff := generateMsgHeader(msgID, msgProps, sequence, phone)
	buff = append(buff, utils.IntTo2byte(packageOption.totalNum)...)
	buff = append(buff, utils.IntTo2byte(packageOption.packetSeq)...)
	return buff
}

func generateCheckSum(buff []byte) byte {
	ret := byte(0)
	for i := 0; i < len(buff); i++ {
		ret = ret ^ buff[i]
	}
	return ret
}

func currSequence() int {
	return 0
}

// 接收消息时 转义
func doEscape4Receive(buff []byte) []byte {
	ret := bytes.Buffer{}
	ret.WriteByte(0x7e)
	for i := 0; i < len(buff)-1; i++ {
		if buff[i] == 0x7d && buff[i+1] == 0x01 {
			ret.WriteByte(0x7d)
			i++
		} else if buff[i] == 0x7d && buff[i+1] == 0x02 {
			ret.WriteByte(0x7e)
			i++
		} else {
			ret.WriteByte(buff[i])
		}
	}
	ret.WriteByte(0x7e)
	return ret.Bytes()
}

// 发送消息时 转义
func doEscape4Send(buff []byte) []byte {
	ret := bytes.Buffer{}
	ret.WriteByte(0x7e)
	for i := 0; i < len(buff); i++ {
		if buff[i] == 0x7e {
			ret.WriteByte(0x7d)
			ret.WriteByte(0x02)
		} else if buff[i] == 0x7d {
			ret.WriteByte(0x7d)
			ret.WriteByte(0x01)
		} else {
			ret.WriteByte(buff[i])
		}
	}
	ret.WriteByte(0x7e)
	return ret.Bytes()
}
func leftPadding(length int, str string) string {
	if length <= len(str) {
		return str[len(str)-length:]
	}
	for i := 0; i < len(str)-length; i++ {
		str = "0" + str
	}
	return str
}
