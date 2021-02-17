package game

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/kcpnet"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"time"
)

//游戏房间类，对应一局游戏
type GameRoom struct {
	ID         int64
	addr       string
	server     *kcpnet.KcpServer
	sessions   map[interface{}]*framework.BaseSession
	dispatcher event.EventDispatcher
	Heros      sync.Map
}

//数据持有；连接者指针列表
func NewGameRoom(address string) *GameRoom {
	s, err := kcpnet.NewKcpServer(address)
	if err != nil {
		return nil
	}
	return &GameRoom{
		ID:         tools.UUID_UTIL.GenerateInt64UUID(),
		addr:       address,
		sessions:   make(map[interface{}]*framework.BaseSession),
		server:     s,
		dispatcher: framework.BaseEventDispatcher{},
		//Heros: make(map[int32]*model.Hero),
	}
}

func (g *GameRoom) GetRoomID() int64 {
	return g.ID
}

func (g *GameRoom) GetHeros() sync.Map{
	return g.Heros
}

//注册连接者
func (g *GameRoom) RegisterConnector(c *framework.BaseSession) error {
	g.sessions[c.Id] = c
	return nil
}

func (g *GameRoom) FetchConnector(sessionId int32) *framework.BaseSession {
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
	if nil == session {
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
		conn.SetWindowSize(4800, 4800)
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

		pbMsg := &pb.GMessage{}
		proto.Unmarshal(buf, pbMsg)
		msg := event2.GMessage{}
		msg.SetRoomId(g.ID)
		m := msg.CopyFromMessage(pbMsg)
		m.SetRoomId(g.ID)
		println(m)
		//若为进入世界业务，则不走消息分发，直接创建会话绑定到玩家ID
		if m.GetCode() == int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY) ||
			m.GetCode() == int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST) {
			g.onEnterGame(m.(*event2.GMessage), session)
			continue
		}
		g.dispatcher.FireEvent(m)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (g *GameRoom) onEnterGame(e *event2.GMessage, s *framework.BaseSession) {
	enterGameNotify := e.Data.(*request.EnterGameRequest)
	s.Id = enterGameNotify.PlayerID
	// todo:先不检测会话存在，放开测试，后期加上
	//if nil==g.FetchConnector(s.Id) {
	//注册会话绑定到玩家id
	g.RegisterConnector(s)
	//}
	//初始化hero加入到对局中
	hero := model.NewHero()
	g.RegisterHero(hero)
	//回包
	data := pb.EnterGameResponse{
		ChangeResult: true,
		HeroId:       hero.ID,
	}
	resp := pb.Response{
		EnterGameResponse: &data,
	}
	msg := pb.GMessage{
		MsgType:  pb.MSG_TYPE_RESPONSE,
		MsgCode:  pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE,
		Response: &resp,
	}
	out, _ := proto.Marshal(&msg)
	GAME_ROOM_MANAGER.Unicast(g.ID, s.Id, out)
}

func (g *GameRoom) RegisterHero(h *model.Hero) {
	hero, _ := g.Heros.Load(h.ID)
	if nil == hero {
		g.Heros.Store(h.ID, h)
	}
}

func (g *GameRoom) ModifyHero(h *model.Hero) {
	g.Heros.Delete(h.ID)
	g.Heros.Store(h.ID, h)
}

func (g *GameRoom) FetchHeros() []*model.Hero {
	heros := make([]*model.Hero, 0)
	//for k, h := range g.Heros {
	//	heros = append(heros, h)
	//}
	g.Heros.Range(func(k, v interface{}) bool {
		heros = append(heros, v.(*model.Hero))
		return true
	})
	return heros
}

//更新英雄位置
func (g *GameRoom) UpdateHeroPos() {
	g.Heros.Range(func(k, v interface{}) bool {
		hero := v.(*model.Hero)
		nowTime := time.Now().UnixNano() / 1e6
		timeElapse := nowTime - hero.UpdateTime
		if timeElapse > int64(time.Second) {
			timeElapse = int64(time.Second)
		}
		hero.UpdateTime = nowTime
		distance := float64(timeElapse) * float64(hero.Speed) / 1000
		x, y := tools.CalXY(distance, hero.HeroDirection)
		hero.HeroPosition.X += x
		hero.HeroPosition.Y += y
		g.Heros.Store(k.(int32), hero)
		return true
	})
}
