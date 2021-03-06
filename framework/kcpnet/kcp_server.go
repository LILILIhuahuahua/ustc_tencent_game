package kcpnet

import (
	"github.com/xtaci/kcp-go"
	"sync"
)

type KcpServer struct {
	mu     sync.Mutex
	addr   string
	Listen *kcp.Listener
	Sess   *kcp.UDPSession
	//broader *framework.Broadcaster
}

// NewKcpServer return a *KcpServer
func NewKcpServer(addr string) (s *KcpServer, err error) {
	ts := new(KcpServer)
	ts.addr = addr
	ts.Listen, err = kcp.ListenWithOptions(addr, nil, 0, 0)
	if err != nil {
		return nil, err
	}
	return ts, err
}

//
//func (s *KcpServer) Run() error {
//	for {
//		conn, err := s.Listen.AcceptKCP()
//		if err != nil {
//			return err
//		}
//		connector := framework.NewConnector(s.broader,*conn)
//		err = s.broader.RegisterConnector(connector)
//		if err != nil {
//			return err
//		}
//		go s.Handle(conn)
//	}
//}
//
//func (s *KcpServer) Handle(conn *kcpnet.UDPSession) {
//	buf := make([]byte, 4096)
//	for {
//		n, err := conn.Read(buf)
//		// 处理buf，改为事件驱动型...
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		err = s.broader.BroadcastAll(buf[:n])
//		//n, err = conn.Write(buf[:n])
//		if err != nil {
//			log.Println(err)
//			return
//		}
//	}
//}
