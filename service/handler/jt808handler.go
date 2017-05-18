package handler

import (
	. "Jt808/tcp"
	//	"encoding/hex"
	"encoding/hex"
	"fmt"
	"jt808"
	"time"
)

func HandleJT808Msg(session *Session) {
	//发送注册包
	obj, ok := session.Get("device")
	if !ok {
		return
	}
	jt808Handler := obj.(jt808.TerminalInfo)
	dev_reg := jt808Handler.GenTerminalRegisterPacket()
	conn := session.Conn
	// 发送注册包
	for {
		conn.Write(dev_reg)
		//等待响应
		buff, err := session.ReadOnePacketWithFixDecoder(&jt808Handler, 10)
		if err != nil {
			fmt.Println(err, " ", buff)
			continue
		}
		// 获得注册包

		jt808Handler.Decode(buff)
		fmt.Println(jt808Handler)
		if jt808Handler.GetAuthCode() != "" {
			fmt.Println(jt808Handler.GetAuthCode())
			break
		}
		fmt.Println(jt808Handler)

	}
	// 发送鉴权

	////等待响应
	for {
		dev_auth := jt808Handler.GenTerminalAuth()
		fmt.Println("dev_auth", hex.EncodeToString(dev_auth))
		conn.Write(dev_auth)
		buff, err := session.ReadOnePacketWithFixDecoder(&jt808Handler, 10)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("resp : ", hex.EncodeToString(buff))
		responseCode, isOk := jt808Handler.GetResponse(buff)
		if isOk && responseCode == 0 {
			break
		} else {
			continue
		}

	}

	// 发送心跳和gps
	//d := jt808Handler.GetHeartBeatInterval()
	timeChan := time.NewTicker(time.Second * 180).C
	for {
		// 发送心跳
		dev_heartbeat := jt808Handler.GenTerminalHeartbeat()
		conn.Write(dev_heartbeat)
		<-timeChan
	}

}
