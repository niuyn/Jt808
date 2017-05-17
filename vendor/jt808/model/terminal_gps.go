package model

import (
	"Jt808/utils"
	"bytes"
)

type TerminalGpsBody struct {
	warn_mark int
	status    int
	lat       float64
	lon       float64
	high      int
	speed     int
	heading   int
	timestr   string
}

func (t *TerminalGpsBody) SetLat(lat float64) *TerminalGpsBody {
	t.lat = lat
	return t
}
func (t *TerminalGpsBody) SetLon(lon float64) *TerminalGpsBody {
	t.lon = lon
	return t
}
func (t *TerminalGpsBody) SetWarnMark(warn int) *TerminalGpsBody {
	t.warn_mark = warn
	return t
}
func (t *TerminalGpsBody) SetStatus(status int) *TerminalGpsBody {
	t.status = status
	return t
}
func (t *TerminalGpsBody) SetHigh(hight int) *TerminalGpsBody {
	t.high = hight
	return t
}
func (t *TerminalGpsBody) SetHeading(heading int) *TerminalGpsBody {
	t.heading = heading
	return t
}
func (t *TerminalGpsBody) SetSpeed(speed int) *TerminalGpsBody {
	t.speed = speed
	return t
}
func (t *TerminalGpsBody) SetTime(timestr string) *TerminalGpsBody {
	t.timestr = timestr
	return t
}

func (t *TerminalGpsBody) GetBytes() []byte {
	var buff = bytes.Buffer{}
	if len(t.timestr) < 12 {
		return buff.Bytes()
	}
	buff.Write(utils.IntTo4byte(t.warn_mark))
	buff.Write(utils.IntTo4byte(t.status))
	buff.Write(utils.IntTo4byte(int(t.lat * 1000000)))
	buff.Write(utils.IntTo4byte(int(t.lon * 1000000)))
	buff.Write(utils.IntTo2byte(t.high))
	buff.Write(utils.IntTo2byte(t.speed))
	buff.Write(utils.IntTo2byte(t.heading))
	buff.Write(utils.EncodeCBCDFromString(t.timestr))
	return buff.Bytes()

}
