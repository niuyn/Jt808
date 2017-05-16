package tcp

import (
	"net"
)

type Session struct {
	Conn     *net.TCPConn
	Keys     map[string]interface{}
	buff     ByteBuff
	Handlers []HandlerFunc
	//service  []HandlerFunc
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
	if c.Keys != nil {
		value, exists = c.Keys[key]
	}
	return
}
func (s *Session) Use(hander ...HandlerFunc) {
	s.Handlers = append(s.Handlers, hander...)
}
func (s *Session) Next() {
	if len(s.Handlers) == 0 {
		return
	}
	for i := 0; i < len(s.Handlers); i++ {
		s.Handlers[i](s.Conn)
	}
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

func (s *Session) readWithProcol(handlerFunc HandlerFunc, timeout int) []byte {
	return nil
}

func (s *Session) read(timeout int) []byte {
	return nil
}
