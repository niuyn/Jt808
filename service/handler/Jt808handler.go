package handler

import (
	"JT808"
	. "Jt808/tcp"
	"time"
)

func HandleJT808Msg(session *Session) {
	//发送注册包
	obj, ok := session.Get("device")
	if !ok {
		return
	}
	dev := obj.(JT808.TerminalInfo)
	dev_reg := dev.GenTerminalRegisterPacket()
	conn := session.Conn
	// 发送注册包
	for {
		conn.Write(dev_reg)
		//等待响应
		conn.SetReadDeadline(time.Now().Add(time.Duration(5 * time.Second)))

		//n, err = conn.Read()

	}
	//
}
