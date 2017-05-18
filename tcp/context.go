package tcp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type MsgDecoderRet int

const (
	msgDecoderRet_NONE MsgDecoderRet = iota
	MsgDecoderRet_OK
	MsgDecoderRet_NOT_OK
	MsgDecoderRet_NEED_DATA
)

type Session struct {
	SimpleTcp  *SimpleTcp
	Conn       *net.TCPConn
	Keys       map[string]interface{}
	buff       *ByteBuff
	Handlers   []HandlerFunc
	IoService  IoServiceFunc
	CurHandler *DecodableHandler
	sendQuene  ReqInfo
	recevQuene RespInfo
}
type DecodableHandler struct {
	hander      HandlerFunc
	decoder_ret MsgDecoderRet
	index       int
}
type RespInfo struct {
	respChan     chan []byte
	respChanSize int
}
type ReqInfo struct {
	reqChan     chan []byte
	reqChanSize int
}
type ByteBuff struct {
	position int
	capacity int
	buff     []byte
}

func (s *Session) Set(key string, value interface{}) {
	if s.Keys == nil {
		s.Keys = make(map[string]interface{})
	}
	s.Keys[key] = value
}
func (s *Session) Get(key string) (value interface{}, exists bool) {
	if s.Keys != nil {
		value, exists = s.Keys[key]
	}
	return
}
func (s *Session) SetReq(size int) *Session {
	s.sendQuene.reqChanSize = size
	if s.sendQuene.reqChan == nil {
		s.sendQuene.reqChan = make(chan []byte, size)
	}
	return s
}
func (s *Session) SetResp(size int) *Session {
	s.recevQuene.respChanSize = size
	if s.recevQuene.respChan == nil {
		s.recevQuene.respChan = make(chan []byte, size)
	}
	return s
}
func (s *Session) Use(hander ...HandlerFunc) {
	if s.CurHandler == nil {
		s.CurHandler = &DecodableHandler{hander[0], MsgDecoderRet_NONE}
	}
	s.Handlers = append(s.Handlers, hander...)
}

func (s *Session) ClearBuff() {
	s.buff.position = 0
}
func (s *Session) PullRemainingBuff() []byte {
	length := s.buff.position
	ret := []byte{}
	ret = append(ret, s.buff.buff[0:length]...)
	s.buff.position = 0
	return ret

}

func (s *Session) GetBuff() []byte {
	ret := []byte{}
	ret = append(ret, s.buff.buff[0:s.buff.position]...)
	return ret
}

func (s *Session) GetBuffWithLength(n int) []byte {
	length := s.buff.position
	ret := []byte{}
	ret = append(ret, s.buff.buff[0:length]...)
	if n < length {
		s.buff.position = length - n
	} else {
		s.buff.position = 0
	}
	return ret

}

// 原生读写
func (s *Session) ReadOnePacketWithFixDecoder(timeoutsec int) ([]byte, error) {
	if s.buff.buff == nil {
		s.buff.buff = make([]byte, s.buff.capacity)
	}
	buffer := s.buff
	conn := s.Conn
	var ret []byte
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutsec) * time.Second))
		n, err := conn.Read(buffer.buff[buffer.position:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			log.Println("empty packet, continue")
			continue
		}
		//fmt.Println(buffer.buff[:n])
		buffer.position += n
		_, ret_, isWhole := s.CurHandler.hander.WholePacket(s)
		if isWhole == MsgDecoderRet_NEED_DATA {
			continue
		} else {
			//ret := s.GetBuffWithLength(length)
			buffer.position = 0
			ret = ret_
			break
		}

	}
	return ret, nil
}

//  非阻塞读
func (s *Session) Read() error {
	if s.buff.buff == nil {
		s.buff.buff = make([]byte, s.buff.capacity)
	}
	buffer := s.buff
	conn := s.Conn
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(s.SimpleTcp.ConnectTimeOut) * time.Second))
		n, err := conn.Read(buffer.buff[buffer.position:])
		if err != nil {
			return err
		}
		if n == 0 {
			log.Println("empty packet, continue")
			continue
		}
		buffer.position += n
		// 选择一正确的协议来处理
		if s.CurHandler.decoder_ret == MsgDecoderRet_OK {
			n, ret_, isWhole := s.CurHandler.hander.WholePacket(s)
			if MsgDecoderRet_NEED_DATA == isWhole {
				continue
			} else {
				//ret := s.GetBuffWithLength(length)
				buffer.position -= n
				select {
				case s.recevQuene.respChan <- ret_:
				default:
					log.Println("recevice quenue is overflow throw a paket :", hex.EncodeToString(<-s.recevQuene.respChan))
				}
			}

		} else if s.CurHandler.decoder_ret == MsgDecoderRet_NOT_OK {
			if s.CurHandler.index >= len(s.Handlers)-1 {
				log.Printf("can`t resove the message  : ", hex.EncodeToString(buffer.buff[:buffer.position]))
				buffer.position = 0
			} else {
				for i := s.CurHandler.index; i < len(s.Handlers); i++ {
					ret := s.Handlers[i].Decodable(s)
					s.CurHandler.hander = s.Handlers[i]
					s.CurHandler.decoder_ret = ret
					s.CurHandler.index = i
					if ret != MsgDecoderRet_NOT_OK {
						break
					}

				}

			}

			// 关链接
		} else if s.CurHandler.decoder_ret == MsgDecoderRet_NEED_DATA {
			s.CurHandler.decoder_ret = s.CurHandler.hander.Decodable(s)

		} else {
			for i := 0; i < len(s.Handlers); i++ {
				ret := s.Handlers[i].Decodable(s)
				if ret != MsgDecoderRet_NOT_OK {
					s.CurHandler.hander = s.Handlers[i]
					s.CurHandler.decoder_ret = ret
					s.CurHandler.index = i
				}
			}
			if s.CurHandler.decoder_ret == msgDecoderRet_NONE {
				log.Printf("can`t resove the message  : ", hex.EncodeToString(buffer.buff[:buffer.position]))
				buffer.position = 0
			}

		}

	}
}

// 阻塞读
func (s *Session) MsutGet(timoutsec int) ([]byte, error) {
	select {
	case ret := <-s.recevQuene.respChan:
		return ret, nil
	case <-time.After(time.Duration(timoutsec) * time.Second):
		log.Printf("get timeout")
	}
	return nil, errors.New(fmt.Sprintf("timeout in %d second \n", timoutsec))
}

// 非阻塞写
func (s *Session) Write() error {
	//
	return nil
}
func (s *Session) DirectWrite(buff []byte) (int, error) {
	return s.Conn.Write(buff)
}

//阻塞请求
func (s *Session) ReqWithResponse(req []byte, timeoutsec int) ([]byte, error) {
	return nil, nil
}
