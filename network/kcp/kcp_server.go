package kcp

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/network"
	"github.com/xtaci/kcp-go"
	"log"
	"sync"
)

type KcpServer struct {
	mu   sync.Mutex
	addr string
	listen *kcp.Listener
	sess *kcp.UDPSession
	broader *network.Broadcaster
}

// NewKcpServer return a *KcpServer
func NewKcpServer(addr string) (s *KcpServer, err error) {
	ts := new(KcpServer)
	ts.addr = addr
	ts.listen,err = kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		return nil, err
	}
	return ts, err
}

func (s *KcpServer) Run() error {
	for {
		conn, err := s.listen.AcceptKCP()
		if err != nil {
			return err
		}
		connector := network.NewConnector(s.broader,*conn)
		err = s.broader.RegisterConnector(connector)
		if err != nil {
			return err
		}
		go s.Handle(conn)
	}
}

func (s *KcpServer) Handle(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		// 处理buf，改为事件驱动型...
		if err != nil {
			log.Println(err)
			return
		}
		err = s.broader.NotifyAll(buf[:n])
		//n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
