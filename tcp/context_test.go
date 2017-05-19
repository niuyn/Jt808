package tcp

import (
	"testing"
	"time"
)

var session = &Session{}

func TestSessionMsutGet(t *testing.T) {
	session.recevQuene.respChan = make(chan []byte, 10)
	go func() {
		time.Sleep(5 * time.Second)
		session.recevQuene.respChan <- []byte{0x01, 0x02}
	}()
	b, err := session.MsutGet(10)
	t.Error(b, err)
}
