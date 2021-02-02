package game

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/kcpnet"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
	"log"
)

//游戏房间类，对应一局游戏
type GameRoom struct {
		ID int64
		addr       string
		server     *kcpnet.KcpServer
		sessions   map[interface{}]*framework.BaseSession
		dispatcher event.EventDispatcher
	}

//数据持有；连接者指针列表
func NewGameRoom(address string) *GameRoom {
	s, err := kcpnet.NewKcpServer(address)
	if err != nil {
		return nil
	}
	return &GameRoom{
		ID: tools.UUID_UTIL.GenerateInt64UUID(),
		addr:       address,
		sessions:   make(map[interface{}]*framework.BaseSession),
		server:     s,
		dispatcher: framework.BaseEventDispatcher{},
	}
}

func (g *GameRoom) GetRoomID() int64{
	return g.ID
}

//注册连接者
func (g *GameRoom) RegisterConnector(c *framework.BaseSession) error {
	g.sessions[c.Id] = c
	return nil
}

func (g *GameRoom) FetchConnector(sessionId int64) *framework.BaseSession {
	//g.sessions[c.Id] = c
	return g.sessions[sessionId]
}

//删除连接者
func (g *GameRoom) DeleteConnector(c *framework.BaseSession) error {
	return nil
}

//广播-全体会话
func (g *GameRoom) BroadcastAll(buff []byte) error {
	for _, session := range g.sessions {
		err := session.SendMessage(buff)
		if nil != err {
			println(err)
			return err
		}
	}
	return nil
}

//单播
func (g *GameRoom) Unicast(buff []byte, sessionId int64) error {
	session := g.sessions[sessionId]
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

func (g *GameRoom) Serv() error {
	for {
		conn, err := g.server.Listen.AcceptKCP()
		if err != nil {
			return err
		}
		session := framework.NewBaseSession(conn)
		if err != nil {
			return err
		}
		go g.Handle(session)
	}
}

func (g *GameRoom) Handle(session *framework.BaseSession) {
	buf := make([]byte, 4096)
	for {
		_, err := session.Sess.Read(buf)
		fmt.Println(string(buf))

		pbMsg:=&pb.GMessage{}
		proto.Unmarshal(buf, pbMsg)
		msg := event2.GMessage{}
		msg.SetRoomId(g.ID)
		m := msg.CopyFromMessage(pbMsg)
		m.SetRoomId(g.ID)
		//若为进入世界业务，则不走消息分发，直接创建会话绑定到玩家ID
		if m.GetCode() == int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY) {
			g.onEnterGame(m.(*event2.GMessage), session)
			return
		}
		g.dispatcher.FireEvent(m)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (g *GameRoom)onEnterGame(e *event2.GMessage, s *framework.BaseSession) {
	enterGameNotify := e.Data.(*notify.EnterGameNotify)
	s.Id = enterGameNotify.PlayerID
	if nil==g.FetchConnector(s.Id) {
		g.RegisterConnector(s)
	}
}
