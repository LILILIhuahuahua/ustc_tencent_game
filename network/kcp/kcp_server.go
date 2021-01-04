package kcp

import (
	"github.com/xtaci/kcp-go"
	"log"
	"sync"
)

type kcpServer struct {
	mu   sync.Mutex
	addr string
	listen *kcp.Listener
	sess *kcp.UDPSession
}

// NewKcpServer return a *kcpServer
func NewKcpServer(addr string) (s *kcpServer, err error) {
	ts := new(kcpServer)
	ts.addr = addr
	ts.listen,err = kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		return nil, err
	}
	return ts, err
}

func (s *kcpServer) Run() error {
	for {
		conn, err := s.listen.AcceptKCP()
		if err != nil {
			return err
		}
		go s.Handle(conn)
	}
}

func (s *kcpServer) Handle(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		_, err := conn.Read(buf)
		// 处理buf，改为事件驱动型...
		if err != nil {
			log.Println(err)
			return
		}

		//n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
