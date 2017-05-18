package jt808

import (
	"Jt808/tcp"
	"Jt808/utils"
	//	"encoding/hex"
	//"fmt"
	. "jt808/consts"
)

func (dev *TerminalInfo) Decodable(session *tcp.Session) tcp.MsgDecoderRet {
	// 检查标识符
	buffer := session.GetBuff()
	if buffer[0] == Msg_Terminal_Identifer {
		return tcp.MsgDecoderRet_OK
	}
	return tcp.MsgDecoderRet_NOT_OK
}
func (dev *TerminalInfo) WholePacket(session *tcp.Session) (int, []byte, tcp.MsgDecoderRet) {
	//长度 较验位
	buffer := session.GetBuff()
	if len(buffer) < 13 {
		return 0, nil, tcp.MsgDecoderRet_NEED_DATA
	}
	if buffer[len(buffer)-1] != 0x7e {
		flag := false
		for i := 1; i < len(buffer); i++ {
			if buffer[i] == Msg_Terminal_Identifer {
				buffer = buffer[0 : i+1]
				flag = true
				break
			}
		}
		if !flag {
			return 0, nil, tcp.MsgDecoderRet_NEED_DATA
		}
	}
	bufferEscaped := doEscape4Receive(buffer[1 : len(buffer)-1])
	// 消息体长度
	msgLen := getLength(bufferEscaped[1:13])
	//
	// 消息头长度
	pros_h := buffer[3]
	// 消息头长度
	msgHeaderlen := 12
	if (pros_h & 0x02) == 0x02 {
		msgHeaderlen = msgHeaderlen + 4
	}
	//消息总长度 标识符 + 消息头 + 消息体 + 较验位 + 标识符
	totalLen := 2 + msgHeaderlen + msgLen + 1
	if totalLen == len(bufferEscaped) {
		// 计算较验位
		checksum := generateCheckSum(bufferEscaped[1 : len(bufferEscaped)-2])
		if checksum == bufferEscaped[len(bufferEscaped)-2] {
			return len(buffer), bufferEscaped, tcp.MsgDecoderRet_OK
		}
	}

	return 0, nil, tcp.MsgDecoderRet_NOT_OK
}

func getLength(header []byte) int {
	//长度
	msgProps := (int(header[2:4][0]<<8) | int(header[2:4][1]))
	// [ 0-9 ] 0000,0011,1111,1111(3FF)(消息体长度)
	return msgProps & 0x03FF

}

func getMsgId(header []byte) int {
	return utils.Bytes2ToInt(header[0:2])
}

func getMsgBody(buff []byte) []byte {
	// 消息体长度
	msgLen := getLength(buff[1:13])
	//
	// 消息头长度
	pros_h := buff[3]
	// 消息头长度
	msgHeaderlen := 13
	if (pros_h & 0x02) == 0x02 {
		msgHeaderlen = msgHeaderlen + 4
	}
	return buff[msgHeaderlen : msgHeaderlen+msgLen]

}

func getSequenceInBuff(msgBody []byte) int {
	return utils.Bytes2ToInt(msgBody[0:2])
}

func (dev *TerminalInfo) Decode(buff []byte) tcp.MsgDecoderRet {
	msgId := getMsgId(buff[1:3])
	switch msgId {
	case Cmd_Common_Resp:
	case Cmd_Terminal_Param_Query:
	case Cmd_Terminal_Tegister_Resp:
		body := getMsgBody(buff)
		dev.authcode = utils.GetStringWithGBK(body[3:])
		return tcp.MsgDecoderRet_OK
	case Cmd_Terminal_Param_Settings:

	}
	return tcp.MsgDecoderRet_NOT_OK

}

func (dev *TerminalInfo) GetResponse(buff []byte) (int, bool) {
	msgId := getMsgId(buff[1:3])
	msgBody := getMsgBody(buff)
	if dev.sequnce == getSequenceInBuff(msgBody) {
		switch msgId {
		case Cmd_Common_Resp:
			return int(msgBody[len(msgBody)-1]), true
		case Cmd_Terminal_Param_Query:
		case Cmd_Terminal_Tegister_Resp:
			dev.authcode = utils.GetStringWithGBK(getMsgBody(buff))
			return int(msgBody[2]), true
		case Cmd_Terminal_Param_Settings:

		}

	}
	return -1, false

}
