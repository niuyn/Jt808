package model

import (
	"JT808monitor/utils"
	"bytes"
)

type TerminalRegisterBody struct {
	province      int    // 省域  word
	city          int    // 市县域 word
	manufacturer  string // 制造商 byte[5]
	terminal_type string // 终端型号byte[20]
	terminal_id   string // 终端ID byte[7]
	/**
	 * 车牌颜色(BYTE) 车牌颜色，按照 JT/T415-2006 的 5.4.12 未上牌时，取值为0<br>
	 * 0===未上车牌<br>
	 * 1===蓝色<br>
	 * 2===黄色<br>
	 * 3===黑色<br>
	 * 4===白色<br>
	 * 9===其他
	 */
	car_license_color      byte   // 车牌颜色 byte
	car_license_identifier string // STRING gbk
}

func (t *TerminalRegisterBody) New() *TerminalRegisterBody {
	return &TerminalRegisterBody{}
}
func (t *TerminalRegisterBody) SetProvince(province int) *TerminalRegisterBody {
	t.province = province
	return t
}
func (t *TerminalRegisterBody) SetCity(city int) *TerminalRegisterBody {
	t.city = city
	return t
}
func (t *TerminalRegisterBody) SetManufacturer(manufacturer string) *TerminalRegisterBody {
	t.manufacturer = manufacturer
	return t
}
func (t *TerminalRegisterBody) SetCarLicenseColor(car_license_color byte) *TerminalRegisterBody {
	if car_license_color >= 0 && car_license_color < 5 {
		t.car_license_color = car_license_color
	} else {
		t.car_license_color = 9
	}
	return t
}
func (t *TerminalRegisterBody) SetTerminalId(terminal_id string) *TerminalRegisterBody {
	t.terminal_id = terminal_id
	return t
}

func (t *TerminalRegisterBody) SetCarLicenseIdentifier(car_license_identifier string) *TerminalRegisterBody {
	t.car_license_identifier = car_license_identifier
	return t
}

func (t *TerminalRegisterBody) getBytes() []byte {
	var buff = bytes.Buffer{}
	buff.Write(utils.IntTo2byte(t.province))
	buff.Write(utils.IntTo2byte(t.province))
       // 制造商
	return nil
}
