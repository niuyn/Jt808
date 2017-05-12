package Jtconsts

const (
	msg_terminal_register     uint = 0x0100 //终端注册
	msg_terminal_logout       uint = 0x0003 //终端注销
	msg_terminal_auth         uint = 0x0102 //终端鉴权
	msg_terminal_gps_up       uint = 0x0200 //位置信息汇报
	msg_terminal_gps_batch_up uint = 0x0704 //批量数据上传
)
