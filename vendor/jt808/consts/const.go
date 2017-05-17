package Jtconsts

const (
	Msg_Terminal_Register     int  = 0x0100 //终端注册
	Msg_Terminal_Logout       int  = 0x0003 //终端注销
	Msg_Terminal_Auth         int  = 0x0102 //终端鉴权
	Msg_Terminal_Gps_Up       int  = 0x0200 //位置信息汇报
	Msg_Terminal_Gps_Batch_Up int  = 0x0704 //批量数据上传
	Msg_Terminal_Heartbeat    int  = 0x0002 //批量数据上传
	Msg_Terminal_Identifer    byte = 0x7e   // 标识符

	// 平台通用应答
	Cmd_Common_Resp int = 0x8001
	// 终端注册应答
	Cmd_Terminal_Tegister_Resp int = 0x8100
	// 设置终端参数
	Cmd_Terminal_Param_Settings int = 0X8103
	// 查询终端参数
	Cmd_Terminal_Param_Query int = 0x8104
)
