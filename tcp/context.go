package tcp

import (
	"log"
	"net"
	"time"
)

type Session struct {
	Conn       *net.TCPConn
	Keys       map[string]interface{}
	buff       *ByteBuff
	Handlers   []HandlerFunc
	IoService  IoServiceFunc
	CurHandler HandlerFunc
	OutChan    chan OutPut
	InputChan  chan Input
}
type OutPut interface {
	SaveTo(obj interface{})
}
type Input interface {
	ReceviveFrom(obj interface{})
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
func (s *Session) Use(hander ...HandlerFunc) {
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

func (s *Session) ReadOnePacketWithFixDecoder(handlerFunc HandlerFunc, timeoutsec int) ([]byte, error) {
	if s.buff.buff == nil {
		s.buff.buff = make([]byte, s.buff.capacity)
	}
	buffer := s.buff
	conn := s.Conn
	var ret []byte
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(10 * time.Second)))
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
		_, ret_, isWhole := handlerFunc.WholePacket(s)
		if !isWhole {
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

func (s *Session) read(timeout int) []byte {
	return nil
}
