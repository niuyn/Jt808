package main

import (
	. "Jt808/common"
	"Jt808/service/handler"
	"Jt808/tcp"
	"fmt"
	"github.com/vaughan0/go-ini"
	"io/ioutil"
	"jt808"
	"log"
	"os"
	"strconv"
	"strings"
)

var GlobalCfg, _ = ini.LoadFile("conf/config.ini")

func main() {
	logger := log.Logger{}
	//读配置文件
	server, ok := GlobalCfg.Get("COMMON", "server")
	if !ok {
		fmt.Println("server ip/port not set")
		return
	}

	_startImei, ok := GlobalCfg.Get("COMMON", "start_imei")
	if !ok {
		fmt.Println("start imei not set")
		return
	}
	points, err := initRoute("conf/route.txt")
	if err != nil {
		fmt.Printf("routes is not ok ")
		return
	}
	simpleTcpClinet := tcp.New(server, 1000)
	// 生成一个模拟器
	terminal := jt808.TerminalInfo{}
	terminal.SetGpsInterval(10).SetHeartBeatInterval(180).SetImei(_startImei)
	err = simpleTcpClinet.HandleService(handler.HandleJT808Msg).Set("device", terminal).Set("routes", &points).Use(&terminal).SetTimeout(1000).Dial()
	if err != nil {
		logger.Println(err)
		return
	}
	simpleTcpClinet.Run()

}
func initRoute(filePath string) ([]Point, error) {
	bytes, err := readAll(filePath)
	if err != nil {
		return nil, err
	}

	route := string(bytes)                         //将读入的字节转换成字符串
	points := strings.Split(route, ";")            //将每组字符串以；标记分开
	pos := make([]Point, len(points), len(points)) //定义一个pos切片
	for index, po := range points {
		values := strings.Split(po, ",") //将每组数据	用，号分开之后存入values，
		//fmt.Println(len(values), index)
		if len(values) != 5 { //判断每组数据的个数，小于5个说明数据不完整,string类型
			fmt.Println("unknown format")
			continue
		}

		pos[index].Lon, _ = strconv.ParseFloat(values[0], 64)
		//pos[index].lat, _ = strconv.ParseFloat(values[1], 64)

		//fmt.Println(pos[index].lon)
		pos[index].Lat, _ = strconv.ParseFloat(values[1], 64)
		pos[index].Speed = 0

	}
	return pos, err
}
func readAll(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}
