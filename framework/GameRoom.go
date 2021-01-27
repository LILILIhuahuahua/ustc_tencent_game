package framework

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/kcpnet"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/golang/protobuf/proto"
	"github.com/xtaci/kcp-go"
	"log"
)

//游戏房间类，对应一局游戏
type GameRoom struct {
		addr     string
		server   *kcpnet.KcpServer
		sessions map[interface{}]*BaseSession
		dispatcher event.EventDispatcher
	}

//数据持有；连接者指针列表
func NewGameRoom(address string) *GameRoom {
	s, err := kcpnet.NewKcpServer(address)
	if err != nil {
		return nil
	}
	return &GameRoom{
		addr:     address,
		sessions: make(map[interface{}]*BaseSession),
		server:   s,
		dispatcher: BaseEventDispatcher{},
	}
}

//注册连接者
func (b *GameRoom) RegisterConnector(c *BaseSession) error {
	b.sessions[c.Id] = c
	return nil
}

//删除连接者
func (b *GameRoom) DeleteConnector(c *BaseSession) error {
	return nil
}

//广播-全体会话
func (b *GameRoom) BroadcastAll(buff []byte) error {
	for _, session := range b.sessions {
		err := session.SendMessage(buff)
		if nil != err {
			println(err)
			return err
		}
	}
	return nil
}

//单播
func (b *GameRoom) Unicast(buff []byte, sessionId int64) error {
	session := b.sessions[sessionId]
	if nil == session{
		return nil
	}
	err := session.SendMessage(buff)
	if nil != err {
		println(err)
		return err
	}
	return nil
}

func (b *GameRoom) Serv() error {
	for {
		conn, err := b.server.Listen.AcceptKCP()
		if err != nil {
			return err
		}
		connector := NewBaseSession(conn)
		err = b.RegisterConnector(connector)
		if err != nil {
			return err
		}
		go b.Handle(connector.Sess)
	}
}

func (b *GameRoom) Handle(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		_, err := conn.Read(buf)
		fmt.Println(string(buf))

		pbMsg:=&pb.GMessage{}
		proto.Unmarshal(buf, pbMsg)
		msg := event2.GMessage{}
		m := msg.CopyFromMessage(pbMsg)
		b.dispatcher.FireEvent(m)
		//req :=&pb.EntityInfoChangeRequest{}
		// 处理buf，改为事件驱动型...
		if err != nil {
			log.Println(err)
			return
		}
		//c.Sess.Write(buf[:n])
		//err = b.BroadcastAll(buf[:n])
		//n, err = conn.Write(buf[:n])
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
	}
}
