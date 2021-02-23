package framework

import (
	"errors"
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
		LastUpdateTime     int64           //上一次接收到消息的时间
		LastDisconnectTime int64           //会话上一次断开时间
		OfflineForever     bool            //超过30s没有发送消息即认为该player永久掉线了
		StatusMutex        sync.Mutex
	}
)

func NewBaseSession(s *kcp.UDPSession) *BaseSession {
	//kcp session调优
	s.SetNoDelay(1, 10, 2, 1)
	s.SetACKNoDelay(true)
	createTime := time.Now().UnixNano() / 1e6
	baseSession := &BaseSession{
		Id:             tools.UUID_UTIL.GenerateInt32UUID(),
		Sess:           s,
		Status:         configs.SessionStatusCreated,
		CreationTime:   createTime,
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

//超过5秒没有收到消息即可认为session挂了,有可能是网络延迟大，也有可能是玩家退出游戏，但此时只是将session的status变为dead，没有将offlineForever变为true
//在此期间如果玩家断线重连的话是可以继续玩游戏的
func (c *BaseSession) IsAvailable() bool {
	nowTime := time.Now().UnixNano() / 1e6
	if c.Status != configs.SessionStatusDead &&
		nowTime-c.LastUpdateTime >= 5*int64(time.Second)/1e6 {
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

func (c *BaseSession) ChangOfflineStatus(status bool) {
	c.OfflineForever = status
}

func (c *BaseSession) CloseKcpSession() error {
	err := c.Sess.Close()
	c.LastDisconnectTime = time.Now().UnixNano() / 1e6
	return err
}

// 30s内如果没有收到玩家的消息，就可以认为玩家永远掉线，需要将offlineForever字段变为true，玩家即使再重连，也需要重新开始
func (c *BaseSession) IsDeprecated() bool {
	nowTime := time.Now().UnixNano() / 1e6
	if !c.OfflineForever &&
		c.Status == configs.SessionStatusDead &&
		nowTime-c.LastDisconnectTime >= 30*int64(time.Second)/1e6 {
		return true
	}
	return false
}
