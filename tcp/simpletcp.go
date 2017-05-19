package tcp

import (
	"errors"
	"log"
	"net"
	"time"
)

type SimpleTcp struct {
	Address        *net.TCPAddr
	ConnectTimeOut int
	session        *Session
	buffSize       int
	serviceStart   bool
}

//type HandlerFunc func(*net.TCPConn)
type IoServiceFunc func(*Session) bool
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
	simpleTcp := &SimpleTcp{tcpAddr, 1000, session, size, false}
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
	if conn.session.sendQuene.reqChanSize == 0 {
		//默认值300
		conn.session.SetReq(300)
	}
	if conn.session.recevQuene.respChanSize == 0 {
		conn.session.SetResp(300)
	}
	if conn.session.IoService == nil {
		return errors.New("sevice can`t be null")
	}
	if conn.session.Handlers == nil {
		return errors.New("handle can`t be null")
	}

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
	// 读写
	if conn.serviceStart {
		log.Printf("sevice is start")
		return
	} else {
		conn.serviceStart = true
		conn.session.IoService(conn.session)

	}
}

func (conn *SimpleTcp) CloseSevice() {
	conn.serviceStart = false
}

func (conn *SimpleTcp) LongRun() {
	cReconnect := make(chan string)
	for {
		go conn.Run()

		go func() {
			err := conn.session.Read()
			if err != nil {
				log.Printf(err.Error())
				cReconnect <- "reconn"
				return
			}

		}()
		go func() {
			err := conn.session.Write(cReconnect)
			if err != nil {
				log.Printf(err.Error())
				return
			}

		}()
		<-cReconnect
		time.Sleep(5 * time.Second)
		conn.Reconn()
	}
}

func (conn *SimpleTcp) Reconn() error {
	return conn.Dial()
}

func StartWebUI() error {
	return nil

}
