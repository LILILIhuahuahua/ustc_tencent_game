package network

import (
	"github.com/google/uuid"
	"github.com/xtaci/kcp-go"
)

//网络组件----连接者
type (
	Connector struct {
		Id string
		//Sess Session
		Sess *kcp.UDPSession
	}
)

func NewConnector(s *kcp.UDPSession) *Connector {
	return &Connector{
		Id:   uuid.New().String(),
		Sess: s,
	}
}

//持有客户端udp信息、session会话

//状态更新
func (c *Connector) Update(buff []byte) error {
	_, err := c.Sess.Write(buff)
	return err
}

//消息接受

//消息移交
