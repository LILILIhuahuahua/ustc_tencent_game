package framework

import (
	"github.com/google/uuid"
	"github.com/xtaci/kcp-go"
)

//基础会话类
type (
	BaseSession struct {
		Id int64				//唯一标识号，与player的ID相同
		Sess *kcp.UDPSession	//kcp发送方
		Status int32			//会话状态：建立、销毁
		Type int32				//网络类型：TCP、UDP
		CreationTime int64		//会话创建时间
		LastDisconnectTime int64//会话上一次断开时间
	}
)

func NewBaseSession(s *kcp.UDPSession) *BaseSession {
	return &BaseSession{
		Id:   int64(uuid.New().ID()),
		Sess: s,
	}
}

//状态更新
func (c *BaseSession) SendMessage(buff []byte) error {
	_, err := c.Sess.Write(buff)
	return err
}

