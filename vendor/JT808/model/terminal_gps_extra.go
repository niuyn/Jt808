package model

import "bytes"

type TerminalGpsExtraBody struct {
	extraMsgID byte
	msgBody    []byte
}

func (t *TerminalGpsExtraBody) SetMsgID(id byte) *TerminalGpsExtraBody {
	t.extraMsgID = id
	return t
}
func (t *TerminalGpsExtraBody) SetMsgBody(body []byte) *TerminalGpsExtraBody {
	t.msgBody = body
	return t
}
func (t *TerminalGpsExtraBody) GetBytes() []byte {
	var buff = bytes.Buffer{}
	buff.Write([]byte{t.extraMsgID, byte(len(t.msgBody))})
	buff.Write(t.msgBody)
	return buff.Bytes()
}
