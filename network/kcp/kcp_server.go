package kcp
import (
	"github.com/xtaci/kcp-go"
	"net"
	"sync"
)

type kcpServer struct {
	mu        sync.Mutex
	addr      string
	ln        net.Listener
}

// NewKcpServer return a *kcpServer
func NewKcpServer(addr string) (s *kcpServer, err error) {
	ts := new(kcpServer)
	ts.addr = addr
	ts.ln, err = kcp.ListenWithOptions(addr, nil, 10, 3)
	if err != nil {
		return nil, err
	}
	return ts, err
}

func (s *kcpServer) Run() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.Handle(conn)
	}
}

func (s *kcpServer) Handle(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
}