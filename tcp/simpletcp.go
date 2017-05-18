package tcp

import (
	"net"
)

type SimpleTcp struct {
	Address        *net.TCPAddr
	ConnectTimeOut int
	session        *Session
	buffSize       int
}

//type HandlerFunc func(*net.TCPConn)
type IoServiceFunc func(*Session)
type HandlerFunc interface {
	Decodable(session *Session) MsgDecoderRet
	WholePacket(session *Session) (int, []byte, MsgDecoderRet)
	Decode(buff []byte) MsgDecoderRet
}

func New(addr string, size int) *SimpleTcp {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil
	}
	session := &Session{}
	simpleTcp := &SimpleTcp{tcpAddr, 1000, session, size}
	session.SimpleTcp = simpleTcp
	return simpleTcp
}
func (conn *SimpleTcp) Use(handler ...HandlerFunc) *SimpleTcp {
	conn.session.Use(handler...)
	return conn
}
func (conn *SimpleTcp) HandleService(serviceFunc IoServiceFunc) *SimpleTcp {
	conn.session.IoService = serviceFunc
	return conn
}
func (conn *SimpleTcp) Set(key string, value interface{}) *SimpleTcp {
	conn.session.Set(key, value)
	return conn
}
func (conn *SimpleTcp) Dial() error {
	realConn, err := net.DialTCP("tcp", nil, conn.Address)
	if err != nil {
		return err
	}
	conn.session.Conn = realConn
	if conn.session.buff == nil {
		conn.session.buff = &ByteBuff{}
	}
	conn.session.buff.capacity = conn.buffSize
	return nil
}

func (conn *SimpleTcp) SetSendBuff(size int) *SimpleTcp {
	conn.session.SetReq(size)
	return conn
}

func (conn *SimpleTcp) SetRecevBuff(size int) *SimpleTcp {
	conn.session.SetResp(size)
	return conn
}
func (conn *SimpleTcp) SetTimeout(timeout int) *SimpleTcp {
	conn.ConnectTimeOut = timeout
	return conn
}
func (conn *SimpleTcp) Run() {

	conn.session.IoService(conn.session)
}

func (conn *SimpleTcp) Reconn() error {
	return conn.Dial()
}

func StartWebUI() error {
	return nil

}
