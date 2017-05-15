package JT808

import (
	. "JT808/consts"
	. "JT808/model"
	"fmt"
)

type TerminalInfo struct {
	imei               string
	heartbeat_interval int
	gps_interval       int
	terminalGps        TerminalGpsInfo
	terminalMultiGps   []TerminalGpsInfo
	terminalReg        TerminalRegisterBody
	sequnce            int
}
type TerminalGpsInfo struct {
	terminalGps      TerminalGpsBody
	terminalGpsExtra TerminalGpsExtraBody
}

// 终端注册
func (dev *TerminalInfo) GenTerminalRegisterPacket() []byte {
	// 生成消息体
	body := dev.terminalReg.GetBytes()
	// 生成消息体属性
	return dev.GenWholeMsg(body, Msg_Terminal_Register)
}

//终端位置信息上报
func (dev *TerminalInfo) GenTerminalGpsUp() []byte {
	//消息体
	gpsBody := dev.terminalGps.terminalGps.GetBytes()
	if len(gpsBody) == 0 {
		fmt.Println("time is empty")
		return ret
	}
	gpsBodyExtra := dev.terminalGps.terminalGpsExtra.GetBytes()
	body := append(gpsBody, gpsBodyExtra...)
	return dev.GenWholeMsg(body, Msg_Terminal_Gps_Up)
}

//终端注销
func (dev *TerminalInfo) GenTerminalLogout() []byte {

	return nil
}

//终端鉴权
func (dev *TerminalInfo) GenTerminalAuth() []byte {
	return nil
}

//定量数据批量上传
func (dev *TerminalInfo) GenTerminalBatchGpsUp() []byte {
	return nil
}

func (dev *TerminalInfo) GenWholeMsg(body []byte, MsgId int) []byte {
	ret := []byte{}
	err, msgProps := generateMsgBodyProps(len(body), 0, 0, false)
	if err != nil {
		fmt.Println(err)
		return ret
	}
	// 生成消息头
	header := generateMsgHeader(MsgId, msgProps, dev.getSequence(), dev.imei)
	ret = append(ret, header...)
	ret = append(ret, body...)
	// 计算校验位
	checksum := generateCheckSum(ret)
	ret = append(ret, checksum)
	return doEscape4Send(ret)
}

func (dev *TerminalInfo) getSequence() int {
	dev.sequnce++
	if dev.sequnce > 0xFFFF {
		dev.sequnce = 0
	}
	return dev.sequnce
}
