package tcp

import "testing"

var session = &Session{}

func TestSessionMsutGet(t *testing.T) {
	session.recevQuene.respChan = make(chan []byte, 10)
	session.recevQuene.respChan <- []byte{0x01, 0x02}
	b, err := session.MsutGet(3)
	t.Error(b, err)
}
