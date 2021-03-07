package game

import (
	"errors"
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/kcpnet"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/aoi"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/prop"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
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
	towers		   []*aoi.Tower
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
		towers:		aoi.InitTowers(),
		//Heros: make(map[int32]*model.Hero),
	}
}

func (g *GameRoom) GetRoomID() int64 {
	return g.ID
}

func (g *GameRoom) GetHeros() sync.Map {
	return g.Heros
}

func (g *GameRoom) GetTowers() []*aoi.Tower {
	return g.towers
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
func (g *GameRoom) Unicast(buff []byte, session *framework.BaseSession) error {
	if session == nil {
		return errors.New("session 为空")
	}
	err := session.SendMessage(buff)
	if nil != err {
		println(err)
		return err
	}
	return nil
}

//多播
func (g *GameRoom) Multiplecast(buff []byte, sessions []*framework.BaseSession) error {
	for _, session := range sessions {
		err := session.SendMessage(buff)
		if nil != err {
			println(err.Error())
			return err
		}
	}
	return nil
}

// 获取某玩家附近的玩家
func (g *GameRoom) GetPlayersNearby(hero *model.Hero) []*model.Hero {
	towers := g.GetTowers()
	var heros []*model.Hero
	var towersOfPlayer []*aoi.Tower
	hero.OtherTowers.Range(func(k, v interface{}) bool {
		towersOfPlayer = append(towersOfPlayer, v.(*aoi.Tower))
		return true
	})
	towersOfPlayer = append(towersOfPlayer, towers[hero.TowerId])
	for _, tower := range towersOfPlayer {
		heros = append(heros, tower.GetHeros()...)
	}
	return heros
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
		//log.Printf("Receive data: %+v", pbMsg)
		msg := event2.GMessage{}
		msg.SetRoomId(g.ID)
		m := msg.CopyFromMessage(pbMsg)
		m.SetRoomId(g.ID)
		//log.Printf("Build event: %+v", m)
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
	hero := model.NewHero(s)
	g.RegisterHero(hero)
	g.SessionHeroMap.Store(s.Id, hero)
	// 视野处理
	towers := g.GetTowers()
	towerId := tools.CalTowerId(hero.HeroDirection.X, hero.HeroDirection.Y)
	otherTowers := tools.GetOtherTowers(towerId)
	hero.TowerId = towerId
	for _, id := range otherTowers {  //存储周围的towerId到hero
		hero.OtherTowers.Store(id, towers[towerId])
	}
	towers[towerId].HeroEnter(hero) //将hero存入tower中
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
		SendTime: tools.TIME_UTIL.NowMillis(),
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

// 在修改hero的位置信息的时候，会将灯塔进行更新
func (g *GameRoom) ModifyHero(modifyHero *model.Hero) {
	heroObj, _ := g.Heros.Load(modifyHero.ID)
	hero := heroObj.(*model.Hero)
	hero.HeroPosition = modifyHero.HeroPosition
	hero.HeroDirection = modifyHero.HeroDirection
	hero.Speed = modifyHero.Speed
	hero.Size = modifyHero.Size
	towers := g.GetTowers()
	towerId := tools.CalTowerId(modifyHero.HeroPosition.X, modifyHero.HeroPosition.Y) // 计算更新位置之后的towerId
	if towerId != hero.TowerId {
		towers[towerId].HeroEnter(hero) // 将hero加入灯塔中
		towers[hero.TowerId].HeroLeave(hero) // 将hero从原来灯塔中删除
		hero.TowerId = towerId
		otherIds := tools.GetOtherTowers(towerId)
		if otherIds == nil {
			fmt.Println("获取其他TowerId的时候出错了")
		}
		midMap := make(map[int32]bool) // 其实相当于一个set，查询元素是否在其中的时间复杂度为O(1)
		for _, id := range otherIds {
			midMap[id] = false
		}
		var needToDelete []*aoi.Tower
		hero.OtherTowers.Range(func(k, v interface{}) bool {
			if _, ok := midMap[k.(int32)]; !ok {  // 如果新的otherTowerId中没有该key，证明该key所对应的tower不在九宫格范围内
				needToDelete = append(needToDelete, v.(*aoi.Tower))
			} else {
				midMap[k.(int32)] = true
			}
			return true
		})
		for _, tower := range needToDelete {
			tower.NotifyHeroMsg(hero, configs.Leave, g.SendHeroViewNotify)
			hero.OtherTowers.Delete(towerId)
		}
		// 接下来处理新加入的otherTowerId
		for k, v := range midMap {
			if !v {
				towers[k].NotifyHeroMsg(hero, configs.Enter, g.SendHeroViewNotify)
				hero.OtherTowers.Store(k, towers[k])
				midMap[k] = true
			}
		}
	}
	g.Heros.Store(hero.ID, hero)
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
	var needToUpdate []*model.Hero
	g.Heros.Range(func(k, v interface{}) bool {
		needToUpdate = append(needToUpdate, v.(*model.Hero))
		return true
	})
	for _, hero := range needToUpdate {
		nowTime := time.Now().UnixNano()
		timeElapse := nowTime - hero.UpdateTime
		hero.UpdateTime = nowTime
		distance := float64(timeElapse) * float64(hero.Speed) / 1e9
		x, y := tools.CalXY(distance, hero.HeroDirection.X, hero.HeroDirection.Y)
		hero.HeroPosition.X += x
		hero.HeroPosition.Y += y
		//需要判断是否出现了碰撞
		g.ModifyHero(hero)
	}
}


