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
func (s *SimpleTcp) ClearBuff() {
	s.session.buff.position = 0
}
func (s *SimpleTcp) PullRemainingBuff() []byte {
	length := s.session.buff.position
	ret := []byte{}
	ret = append(ret, s.session.buff.buff[0:length]...)
	s.session.buff.position = 0
	return ret

}
