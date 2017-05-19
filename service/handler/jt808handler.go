package handler

import (
	. "Jt808/tcp"
	//	"encoding/hex"
	. "Jt808/common"
	"encoding/hex"
	"fmt"
	"jt808"
	. "jt808/model"
	"log"
	"time"
)

func HandleJT808Msg(session *Session) bool {
	// gps 和 心跳任务标志
	flag := false
	// 重连标志
	cReconnect := make(chan string)
	//　阻出写方法
	signalRet := make(chan string)
	for {
		//写 conn
		go func() {
			err := session.Write(signalRet)
			if err != nil {
				log.Printf(err.Error())
				return
			}

		}()
		//读conn 到 chan
		go func() {
			err := session.Read()
			if err != nil {
				log.Printf(err.Error())
				cReconnect <- "reconn"
				signalRet <- "return"
				return
			}

		}()
		///发送注册包
		obj, ok := session.Get("device")
		if !ok {
			session.Close()
			return false
		}
		jt808Handler := obj.(jt808.TerminalInfo)
		dev_reg := jt808Handler.GenTerminalRegisterPacket()
		for {
			fmt.Println(hex.EncodeToString(dev_reg))
			session.DirectWrite(dev_reg)
			buff, err := session.MsutGet(10)
			if err != nil {
				time.Sleep(10 * time.Second)
				continue
			}
			//判读是否成功
			jt808Handler.Decode(buff)
			//fmt.Println(jt808Handler)
			if jt808Handler.GetAuthCode() != "" {
				//fmt.Println(jt808Handler.GetAuthCode())
				break
			} else {
				continue
			}

		}
		for {
			dev_auth := jt808Handler.GenTerminalAuth()
			session.DirectWrite(dev_auth)
			buff, err := session.MsutGet(10)
			if err != nil {
				fmt.Println(err)
				continue
			}
			responseCode, isOk := jt808Handler.GetResponse(buff)
			if isOk && responseCode == 0 {
				break
			} else {
				time.Sleep(5 * time.Second)
				continue
			}

		}

		// 发送心跳和gps
		//d := jt808Handler.GetHeartBeatInterval()
		if !flag {
			timeChanForHeartBeat := time.NewTicker(time.Second * time.Duration(jt808Handler.GetHeartBeatInterval())).C
			go func() {
				for {
					// 发送心跳
					dev_heartbeat := jt808Handler.GenTerminalHeartbeat()
					session.WriteSessionBuffer(dev_heartbeat)
					<-timeChanForHeartBeat
				}

			}()

			timeChanForGps := time.NewTicker(time.Second * time.Duration(jt808Handler.GetGpsInterval())).C
			go func() {
				route, ok := session.Get("routes")
				if !ok {
					log.Printf("routes is not ok")
					return
				}
				points_ := route.(*[]Point)
				points := *points_
				index := 0
				for {
					if index >= len(points) {
						index = 0
					}
					point := points[index]
					// 发送GPS
					gpsBody := &TerminalGpsBody{}
					timestr := time.Now().Format("060102150405")
					gpsBody.SetLat(point.Lat).SetLon(point.Lon).SetSpeed(point.Speed).SetTime(timestr)
					jt808Handler.SetGps(jt808.TerminalGpsInfo{*gpsBody, TerminalGpsExtraBody{}})
					dev_gps := jt808Handler.GenTerminalGpsUp()
					fmt.Println(hex.EncodeToString(dev_gps))
					session.WriteSessionBuffer(dev_gps)
					<-timeChanForGps
				}

			}()
			flag = true
		}

		<-cReconnect
		time.Sleep(5 * time.Second)
		session.SimpleTcp.Reconn()
	}

	// 启动是否需要读数据
	return true

}

func HandleJT808Msg_(session *Session) {
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
		buff, err := session.ReadOnePacketWithFixDecoder(10)
		if err != nil {
			fmt.Println(err, " ", buff)
			time.Sleep(5 * time.Second)
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
		buff, err := session.ReadOnePacketWithFixDecoder(10)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("resp : ", hex.EncodeToString(buff))
		responseCode, isOk := jt808Handler.GetResponse(buff)
		if isOk && responseCode == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
			continue
		}

	}

	// 发送心跳和gps
	//d := jt808Handler.GetHeartBeatInterval()
	timeChanForHeartBeat := time.NewTicker(time.Second * time.Duration(jt808Handler.GetHeartBeatInterval())).C
	go func() {
		for {
			// 发送心跳
			dev_heartbeat := jt808Handler.GenTerminalHeartbeat()
			conn.Write(dev_heartbeat)
			<-timeChanForHeartBeat
		}

	}()

	timeChanForGps := time.NewTicker(time.Second * time.Duration(jt808Handler.GetHeartBeatInterval())).C
	go func() {
		for {
			// 发送心跳
			dev_heartbeat := jt808Handler.GenTerminalGpsUp()
			conn.Write(dev_heartbeat)
			<-timeChanForGps
		}

	}()

}
