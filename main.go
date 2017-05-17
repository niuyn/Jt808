package main

import (
	"Jt808/service/handler"
	"Jt808/tcp"
	"jt808"
	"log"
)

func main() {
	logger := log.Logger{}
	simpleTcp := tcp.New("127.0.0.1:8001", 1000)
	// 生成一个模拟器
	terminal := jt808.TerminalInfo{}
	terminal.SetGpsInterval(10).SetHeartBeatInterval(180).SetImei("014530399199")

	err := simpleTcp.HandleService(handler.HandleJT808Msg).Set("device", terminal).Dial()
	if err != nil {
		logger.Println(err)
		return
	}
	simpleTcp.Run()

}
