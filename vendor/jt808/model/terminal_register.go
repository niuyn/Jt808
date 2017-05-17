package model

import (
	"Jt808/utils"
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

func (t *TerminalRegisterBody) GetBytes() []byte {
	var buff = bytes.Buffer{}
	buff.Write(utils.IntTo2byte(t.province))
	buff.Write(utils.IntTo2byte(t.city))
	// 制造商ID byte[5]
	//manufacturer := utils.EncodeCBCDFromString(t.manufacturer)

	buff.Write(rightPadWithZero(5, []byte(t.manufacturer)))
	// 终端型号 byte[20]
	buff.Write(rightPadWithZero(20, utils.EncodeCBCDFromString(t.terminal_type)))
	// 终端ID
	buff.Write(rightPadWithZero(7, utils.EncodeCBCDFromString(t.terminal_id)))
	// 车牌颜色
	buff.Write([]byte{t.car_license_color})
	// 车牌号
	if t.car_license_identifier == "" {
		t.car_license_identifier = "浙A00000"
	}
	buff.Write(utils.GetBytesWithGBK(t.car_license_identifier))

	return buff.Bytes()
}

func rightPadWithZero(length int, buff []byte) []byte {
	if length <= len(buff) {
		return buff[:length]
	}
	delta := length - len(buff)
	for i := 0; i < delta; i++ {
		buff = append(buff, 0x0)
	}

	return buff
}
