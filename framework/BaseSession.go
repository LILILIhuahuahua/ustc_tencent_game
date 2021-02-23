package framework

import (
	"errors"
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/xtaci/kcp-go"
	"sync"
	"time"
)

//基础会话类
type (
	BaseSession struct {
		Id                 int32           //唯一标识号，与player的ID相同
		Sess               *kcp.UDPSession //kcp发送方
		Status             int32           //会话状态：建立、销毁
		Type               int32           //网络类型：TCP、UDP
		CreationTime       int64           //会话创建时间
		LastUpdateTime	   int64  		   //上一次接收到消息的时间
		LastDisconnectTime int64           //会话上一次断开时间
		StatusMutex		   sync.Mutex
	}
)

func NewBaseSession(s *kcp.UDPSession) *BaseSession {
	//kcp session调优
	s.SetNoDelay(1, 10, 2, 1)
	s.SetACKNoDelay(true)
	createTime := time.Now().UnixNano() / 1e6
	baseSession := &BaseSession{
		Id:   tools.UUID_UTIL.GenerateInt32UUID(),
		Sess: s,
		Status: configs.SessionStatusCreated,
		CreationTime: createTime,
		LastUpdateTime: createTime,
	}
	return baseSession
}

//状态更新
func (c *BaseSession) SendMessage(buff []byte) error {
	if c.Status == configs.SessionStatusDead {
		return errors.New("该session已经关闭，不能写入数据")
	}
	_, err := c.Sess.Write(buff)
	return err
}

func (c *BaseSession) IsAvailable() bool {
	nowTime := time.Now().UnixNano() / 1e6
	fmt.Println(nowTime - c.LastUpdateTime)
	fmt.Println(5 * int64(time.Second) / 1e6)
	if nowTime - c.LastUpdateTime >= 5 * int64(time.Second) / 1e6 && c.Status != configs.SessionStatusDead  { //超过5秒没有收到消息即可认为session挂了
		return false
	}
	return true
}

func (c *BaseSession) UpdateTime() {
	c.LastUpdateTime = time.Now().UnixNano() / 1e6
}

func (c *BaseSession) ChangeStatus(status int32) {
	c.StatusMutex.Lock()
	c.Status = status
	c.StatusMutex.Unlock()
	println("session的status变为了", c.Status)
}

func (c *BaseSession) CloseKcpSession() error {
	err := c.Sess.Close()
	c.LastDisconnectTime = time.Now().UnixNano() / 1e6
	return err
}