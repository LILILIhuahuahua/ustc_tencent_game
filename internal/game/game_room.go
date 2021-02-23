package game

import (
	"errors"
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/kcpnet"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/prop"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"time"
)

//游戏房间类，对应一局游戏
type GameRoom struct {
	ID             int64
	addr           string
	server         *kcpnet.KcpServer
	sessions       sync.Map //map[interface{}]*framework.BaseSession
	dispatcher     event.EventDispatcher
	Heros          sync.Map
	SessionHeroMap sync.Map //map[sessionId] *model.Hero
	props          *prop.PropsManger
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
		server:     s,
		dispatcher: framework.BaseEventDispatcher{},
		props:      prop.New(),
		//Heros: make(map[int32]*model.Hero),
	}
}

func (g *GameRoom) GetRoomID() int64 {
	return g.ID
}

func (g *GameRoom) GetHeros() sync.Map {
	return g.Heros
}

//注册连接者
func (g *GameRoom) RegisterConnector(c *framework.BaseSession) error {
	g.sessions.Store(c.Id, c)
	return nil
}

func (g *GameRoom) FetchConnector(sessionId int32) *framework.BaseSession {
	sess, ok := g.sessions.Load(sessionId)
	if !ok {
		return nil
	}
	return sess.(*framework.BaseSession)
}

//删除连接者
func (g *GameRoom) DeleteConnector(c *framework.BaseSession) error {
	return nil
}

//广播-全体会话
func (g *GameRoom) BroadcastAll(buff []byte) error {
	var sendQueue []*framework.BaseSession
	g.sessions.Range(func(_, v interface{}) bool {
		sess := v.(*framework.BaseSession)
		if sess.Status != configs.SessionStatusDead {
			sendQueue = append(sendQueue, v.(*framework.BaseSession))
		}
		return true
	})
	for _, session := range sendQueue {
		err := session.SendMessage(buff)
		if nil != err {
			println(err.Error())
			return err
		}
	}
	return nil
}

//单播
func (g *GameRoom) Unicast(buff []byte, sessionId int64) error {
	session, ok := g.sessions.Load(sessionId)
	if !ok {
		return nil
	}
	err := session.(*framework.BaseSession).SendMessage(buff)
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
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("在handle的时候发生了错误, 错误为：", err)
		}
	}()

	for {
		session.StatusMutex.Lock()
		if session.Status == configs.SessionStatusDead {
			println("baibai")
			break
		}
		//给session加一个读超时函数
		err := session.Sess.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(2)))
		if err != nil {
			panic("setDeadLine出错")
		}
		num, _ := session.Sess.Read(buf)
		session.StatusMutex.Unlock()
		if num == 0 {
			continue
		}
		session.UpdateTime()

		pbMsg := &pb.GMessage{}
		proto.Unmarshal(buf, pbMsg)
		log.Printf("Receive data: %+v", pbMsg)
		msg := event2.GMessage{}
		msg.SetRoomId(g.ID)
		m := msg.CopyFromMessage(pbMsg)
		m.SetRoomId(g.ID)
		log.Printf("Build event: %+v", m)
		//若为进入世界业务，则不走消息分发，直接创建会话绑定到玩家ID
		if m.GetCode() == int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY) ||
			m.GetCode() == int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST) {
			g.onEnterGame(m.(*event2.GMessage), session)
			continue
		}
		g.dispatcher.FireEvent(m)
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
	g.SessionHeroMap.Store(s.Id, hero)
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
		x, y := tools.CalXY(distance, hero.HeroDirection.X, hero.HeroDirection.Y)
		hero.HeroPosition.X += x
		hero.HeroPosition.Y += y
		g.Heros.Store(k.(int32), hero)
		return true
	})
}

func (g *GameRoom) DeleteUnavailableSession() error {
	var needDelete []*framework.BaseSession
	//将不能正常通信的session存储到needDelete中
	g.sessions.Range(func(_, obj interface{}) bool {
		sess := obj.(*framework.BaseSession)
		if !sess.IsAvailable() {
			needDelete = append(needDelete, sess)
		}
		return true
	})
	//fmt.Println(needDelete)
	for _, session := range needDelete {
		session.ChangeStatus(configs.SessionStatusDead)
		err := session.CloseKcpSession()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GameRoom) DeleteOfflinePlayer() error {
	var needDelete []*framework.BaseSession
	g.sessions.Range(func(_, obj interface{}) bool {
		sess := obj.(*framework.BaseSession)
		if sess.IsDeprecated() {
			needDelete = append(needDelete, sess)
		}
		return true
	})

	for _, session := range needDelete {
		deletedObj, ok := g.SessionHeroMap.Load(session.Id)
		if !ok {
			return errors.New("玩家不存在")
		}
		hero := deletedObj.(*model.Hero)
		hero.ChangeHeroStatus(configs.Dead)
		session.ChangOfflineStatus(true)
		fmt.Println("我调用了玩家删除函数")
	}
	return nil
}
