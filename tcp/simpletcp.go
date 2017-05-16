package tcp

import (
	"net"
)

type SimpleTcp struct {
	Address        *net.TCPAddr
	ConnectTimeOut int
	session        *Session
}
type HandlerFunc func(*net.TCPConn)

func (conn *SimpleTcp) New(addr string) *SimpleTcp {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil
	}
	return &SimpleTcp{tcpAddr, 1000, &Session{}}
}
func (conn *SimpleTcp) Use(handler ...HandlerFunc) *SimpleTcp {
	conn.session.Use(handler...)
	return conn
}
func (conn *SimpleTcp) Dial() error {
	realConn, err := net.DialTCP("tcp", nil, conn.Address)
	if err != nil {
		return err
	}
	conn.session.Conn = realConn
	return nil
}

func (conn *SimpleTcp) run() {
	conn.session.Next()
}

func (conn *SimpleTcp) reconn() error {
	return conn.Dial()
}
