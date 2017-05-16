package JT808

import (
	. "JT808/consts"
	. "JT808/model"
	"Jt808/utils"
	"fmt"
)

type TerminalInfo struct {
	imei               string
	heartbeat_interval int
	gps_interval       int
	terminalGps        TerminalGpsInfo
	terminalBatchGps   TerminalBatchGpsInfo
	terminalReg        TerminalRegisterBody
	sequnce            int
	authcode           string
	reconnflag         bool
}

type TerminalGpsInfo struct {
	terminalGps      TerminalGpsBody
	terminalGpsExtra TerminalGpsExtraBody
}

type TerminalBatchGpsInfo struct {
	terminalBatchGpsInfoType int
	terminalMultiGps         []TerminalGpsInfo
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
		return []byte{}
	}
	gpsBodyExtra := dev.terminalGps.terminalGpsExtra.GetBytes()
	body := append(gpsBody, gpsBodyExtra...)
	return dev.GenWholeMsg(body, Msg_Terminal_Gps_Up)
}

//终端注销
func (dev *TerminalInfo) GenTerminalHeartbeat() []byte {

	return dev.GenWholeMsg([]byte{}, Msg_Terminal_Heartbeat)
}

//终端心跳
func (dev *TerminalInfo) GenTerminalLogout() []byte {

	return dev.GenWholeMsg([]byte{}, Msg_Terminal_Logout)
}

//终端鉴权
func (dev *TerminalInfo) GenTerminalAuth() []byte {
	//消息体
	body := utils.GetBytesWithGBK(dev.authcode)
	return dev.GenWholeMsg(body, Msg_Terminal_Auth)
}

//定量数据批量上传
func (dev *TerminalInfo) GenTerminalBatchGpsUp() []byte {
	// 消息体
	//数据项个数
	body := []byte{}
	gpsNum := len(dev.terminalBatchGps.terminalMultiGps)
	gpsType := dev.terminalBatchGps.terminalBatchGpsInfoType
	body = append(body, utils.IntTo2byte(gpsNum)...)
	body = append(body, byte(gpsType))
	//多个位置数据体
	batchGps := []byte{}
	for i := 0; i < gpsNum; i++ {
		gpsBody := dev.terminalBatchGps.terminalMultiGps[i].terminalGps.GetBytes()
		gpsBodyExtra := dev.terminalBatchGps.terminalMultiGps[i].terminalGpsExtra.GetBytes()
		batchGps = append(batchGps, gpsBody...)
		batchGps = append(batchGps, gpsBodyExtra...)
	}
	//获得位置汇报数据体长度
	length := len(batchGps)
	body = append(body, utils.IntTo2byte(length)...)
	body = append(body, batchGps...)

	return dev.GenWholeMsg(body, Msg_Terminal_Gps_Batch_Up)
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
