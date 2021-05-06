package game

import (
	"errors"
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/aoi"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/collision"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	notify2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/prop"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

//GameRoom 游戏房间类，对应一局游戏
type GameRoom struct {
	ID int64
	//addr             string
	//server           *kcpnet.KcpServer
	acceptedSessions sync.Map
	sessions         sync.Map //map[interface{}]*framework.BaseSession
	dispatcher       event.EventDispatcher
	Heroes           sync.Map
	SessionHeroMap   sync.Map //map[sessionId] *model.Hero
	props            *prop.PropsManger
	towers           []*aoi.Tower
	quadTree         *collision.QuadTree //对局内四叉树，用于进行碰撞检测
	heroRankHeap     *GameRankHeap
	gameOver         int32 // 如果 gameOver 为 0，则表示对局仍在继续，当设置为 1 时，表示对局已经结束，此时回收线程
	AliveHeroNum     int32
}

//数据持有；连接者指针列表
func NewGameRoom() *GameRoom {
	//s, err := kcpnet.NewKcpServer(address)
	//if err != nil {
	//	return nil
	//}
	gameroom := &GameRoom{
		ID: tools.UUID_UTIL.GenerateInt64UUID(),
		//addr:         address,
		//server:       s,
		dispatcher:   framework.NewBaseEventDispatcher(configs.MaxEventQueueSize),
		props:        prop.New(),
		towers:       aoi.InitTowers(),
		quadTree:     collision.NewQuadTree("0", 0, collision.NewRectangleByBounds(configs.MapMinX, configs.MapMinY, configs.MapMaxX, configs.MapMaxY)),
		heroRankHeap: NewGameRankHeap(configs.HeroRankListLength),
		//Heroes: make(map[int32]*model.Hero),
		AliveHeroNum: 0,
	}
	roomInitProps, err := gameroom.props.GetProps()
	if err != nil {
		log.Printf("[GameRoom]在初始化道具视野的时候出错了 \n")
	}
	gameroom.AdjustPropsIntoTower(roomInitProps)
	log.Printf("[GameRoom]初始化新对局！GameRoom：%v \n", gameroom)
	return gameroom
}

func (g *GameRoom) GetRoomID() int64 {
	return g.ID
}

func (g *GameRoom) GetHeros() sync.Map {
	return g.Heroes
}

func (g *GameRoom) GetTowers() []*aoi.Tower {
	return g.towers
}

func (g *GameRoom) GetHero(heroId int32) *model.Hero {
	hero, ok := g.Heroes.Load(heroId)
	if !ok {
		return nil
	}
	return hero.(*model.Hero)
}

//注册连接者
func (g *GameRoom) RegisterConnector(c *framework.BaseSession) error {
	g.sessions.Store(c.Id, c)
	return nil
}

// @title    AcceptConnector
// @description 接收GameRoomManager下方的新会话
func (g *GameRoom) AcceptConnector(session *framework.BaseSession) {
	g.acceptedSessions.Store(session.Id, session)
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

// Serv
// @description   游戏服务主方法（负责处理：1.接受连接创建会话 2.监听已接收但还未注册的会话，接收到进入世界请求时注册会话 3.监听已注册会话，投递网络消息包至消息队列中）
func (g *GameRoom) Serv() error {
	go g.HandleSessions()       //开启会话监听线程，监听session集合中的读事件，将读到的GMessage放入环形队列中
	go g.HandleEventFromQueue() //开启消费线程，从环形队列中读取GMessage消息并处理
	go g.UpdateHeros()          // 更新hero的信息（位置、状态等）
	go g.PeriodicalInitProps()  // 定期生成新的道具

	for g.gameOver == 0 {
		//conn, err := g.server.Listen.AcceptKCP()
		//if err != nil {
		//	return err
		//}
		//conn.SetWindowSize(4800, 4800)
		//session := framework.NewBaseSession(conn)
		//if err != nil {
		//	return err
		//}
		//g.acceptedSessions.Store(session.Id, session) //将新会话放入未注册会话集合中
		g.registerSessions() //处理会话注册流程（等待玩家进入世界enterWorld）
	}

	return nil
}

// @title    registerSessions
// @description 监听已接收但还未注册的会话，接收到进入世界请求时注册会话
func (g *GameRoom) registerSessions() {
	buf := make([]byte, 4096)
	g.acceptedSessions.Range(func(_, v interface{}) bool {
		session := v.(*framework.BaseSession)
		err := session.Sess.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(2)))
		if err != nil {
			panic("setDeadLine出错")
		}
		num, _ := session.Sess.Read(buf)
		if num == 0 {
			return true
		}
		session.UpdateTime()
		pbMsg := &pb.GMessage{}
		proto.Unmarshal(buf, pbMsg)
		//log.Printf("Receive data: %+v", pbMsg)
		msg := event2.GMessage{}
		msg.SetRoomId(g.ID)
		m := msg.CopyFromMessage(pbMsg)
		m.SetRoomId(g.ID)
		if m.GetCode() == int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST) {
			// 将该session移出未注册会话集合
			g.acceptedSessions.Delete(session.Id)
			// 进入世界处理，将session放入已注册会话集合
			g.onEnterGame(m.(*event2.GMessage), session)
		}
		//buf清零
		for i := range buf {
			buf[i] = 0
		}
		return true
	})
}

// @title    registerSessions
// @description 监听已接收但还未注册的会话，接收到进入世界请求时注册会话
func (g *GameRoom) HandleSessions() {
	buf := make([]byte, 4096)
	for g.gameOver == 0 {
		g.sessions.Range(func(_, v interface{}) bool {
			session := v.(*framework.BaseSession)
			//处理单个会话的消息读取
			g.Handle(session, buf)
			return true
		})
	}
}

func (g *GameRoom) Handle(session *framework.BaseSession, buf []byte) {
	//buf := make([]byte, 4096)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("在handle的时候发生了错误, 错误为：", err)
		}
	}()
	session.StatusMutex.Lock()
	if session.Status == configs.SessionStatusDead {
		println("会话状态为死亡，无法从中读取数据")
	}
	//给session加一个读超时函数，及时释放锁
	err := session.Sess.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(2)))
	if err != nil {
		panic("setDeadLine出错")
	}
	num, _ := session.Sess.Read(buf)
	session.StatusMutex.Unlock()
	if num == 0 {
		return
	}
	session.UpdateTime()

	pbMsg := &pb.GMessage{}
	proto.Unmarshal(buf, pbMsg)
	//log.Printf("Receive data: %+v", pbMsg)
	msg := event2.GMessage{}
	msg.SetRoomId(g.ID)
	m := msg.CopyFromMessage(pbMsg)
	m.SetRoomId(g.ID)
	//buf清零
	for i := range buf {
		buf[i] = 0
	}
	//放入消息队列中
	g.dispatcher.FireEvent(m)
}

func (g *GameRoom) HandleEventFromQueue() {
	for g.gameOver == 0 {
		e, err := g.dispatcher.GetEventQueue().Pop()
		if nil == e { //todo
			continue
		}
		if nil != err {
			fmt.Println(err)
			continue
		}
		msg := e.(*event2.GMessage)
		framework.EVENT_HANDLER.OnEvent(msg)
	}
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

// 获取某玩家附近的玩家和道具
func (g *GameRoom) GetItemsNearby(hero *model.Hero) ([]*model.Hero, []*model.Prop) {
	towers := g.GetTowers()
	var heros []*model.Hero
	var props []*model.Prop
	var towersOfPlayer []*aoi.Tower
	hero.OtherTowers.Range(func(k, v interface{}) bool {
		towersOfPlayer = append(towersOfPlayer, v.(*aoi.Tower))
		return true
	})
	towersOfPlayer = append(towersOfPlayer, towers[hero.TowerId])
	for _, tower := range towersOfPlayer {
		heros = append(heros, tower.GetHeros()...)
		props = append(props, tower.GetProps()...)
	}
	return heros, props
}

func (g *GameRoom) onEnterGame(e *event2.GMessage, s *framework.BaseSession) {
	enterGameReq := e.Data.(*request.EnterGameRequest)
	s.Id = enterGameReq.PlayerID
	log.Printf("[GameRoom]玩家进入房间！session：%v, room: %v", s, g)
	// todo:先不检测会话存在，放开测试，后期加上
	//if nil==g.FetchConnector(s.Id) {
	//注册会话绑定到玩家id
	g.RegisterConnector(s)
	//}
	//初始化hero加入到对局中
	hero := model.NewHero(enterGameReq.PlayerName, s)
	g.SessionHeroMap.Store(s.Id, hero)
	atomic.AddInt32(&g.AliveHeroNum, 1)
	// 视野处理
	towers := g.GetTowers()
	towerId := tools.CalTowerId(hero.HeroDirection.X, hero.HeroDirection.Y)
	otherTowers := tools.GetOtherTowers(towerId)
	hero.TowerId = towerId
	for _, id := range otherTowers { //存储周围的towerId到hero
		hero.OtherTowers.Store(id, towers[id])
	}
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
		SeqId:    enterGameReq.SeqId,
		Response: &resp,
		SendTime: tools.TIME_UTIL.NowMillis(),
	}
	out, _ := proto.Marshal(&msg)
	GAME_ROOM_MANAGER.Unicast(g.ID, s.Id, out)
	g.RegisterHero(hero)
	//发出排行榜推送
	rankInfos := g.heroRankHeap.GetSortedHeroRankInfos()
	notify := notify2.NewGameRankListNotify(rankInfos)
	GAME_ROOM_MANAGER.Unicast(g.ID, s.Id, notify.ToGMessageBytes())
	//调整hero的注册位置
	towers[towerId].HeroEnter(hero) //将hero存入tower中
	g.NotifyHeroPropMsg(hero)       // 向该hero发送附近的道具信息
}

func (g *GameRoom) RegisterHero(h *model.Hero) {
	hero, _ := g.Heroes.Load(h.ID)
	if nil == hero {
		g.Heroes.Store(h.ID, h)
		//更新排行榜
		heroRankInfo := info.NewHeroRankInfo(h)
		g.heroRankHeap.ChallengeRank(heroRankInfo)
	}
}

// 在修改hero的位置信息的时候，会将灯塔进行更新
func (g *GameRoom) ModifyHero(modifyHero *model.Hero) {
	hero := g.GetHero(modifyHero.ID)
	hero.HeroPosition = modifyHero.HeroPosition
	hero.HeroDirection = modifyHero.HeroDirection
	hero.Size = modifyHero.Size
	towers := g.GetTowers()
	towerId := tools.CalTowerId(modifyHero.HeroPosition.X, modifyHero.HeroPosition.Y) // 计算更新位置之后的towerId
	if towerId >= int32(len(towers)) || towerId < 0 {
		panic("计算towerId时出错")
	}
	if towerId != hero.TowerId {
		towers[towerId].HeroEnter(hero)      // 将hero加入灯塔中
		g.NotifyHeroPropMsg(hero)            // 向该hero发送附近的道具信息
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
		hero.OtherTowers.Delete(towerId) // 要把当前最新的TowerId在原来的OhterTowers里面删除，不然会报错
		hero.OtherTowers.Range(func(k, v interface{}) bool {
			if _, ok := midMap[k.(int32)]; !ok { // 如果新的otherTowerId中没有该key，证明该key所对应的tower不在九宫格范围内
				needToDelete = append(needToDelete, v.(*aoi.Tower))
			} else {
				midMap[k.(int32)] = true
			}
			return true
		})
		for _, tower := range needToDelete {
			g.NotifyHeroView(hero, configs.Leave, tower)
			hero.OtherTowers.Delete(tower.GetId())
		}
		// 接下来处理新加入的otherTowerId
		for k, v := range midMap {
			if !v {
				g.NotifyHeroView(hero, configs.Enter, towers[k])
				hero.OtherTowers.Store(k, towers[k])
				midMap[k] = true
			}
		}
	}
	g.Heroes.Store(hero.ID, hero)
}

func (g *GameRoom) FetchHeros() []*model.Hero {
	heros := make([]*model.Hero, 0)
	//for k, h := range g.Heroes {
	//	heros = append(heros, h)
	//}
	g.Heroes.Range(func(k, v interface{}) bool {
		heros = append(heros, v.(*model.Hero))
		return true
	})
	return heros
}

func (g *GameRoom) UpdateHeros() {
	for atomic.LoadInt32(&g.gameOver) == 0 {
		g.UpdateHeroPosAndStatus()
		time.Sleep(50 * 1e6) //睡50ms
	}
}

//更新英雄位置
func (g *GameRoom) UpdateHeroPosAndStatus() {
	var needToUpdate []*model.Hero
	g.Heroes.Range(func(k, v interface{}) bool {
		needToUpdate = append(needToUpdate, v.(*model.Hero))
		return true
	})
	for _, hero := range needToUpdate {
		if hero.Status == configs.HeroStatusDead {
			continue
		}
		nowTime := time.Now().UnixNano()
		// 处理玩家的无敌时间
		if hero.Invincible && nowTime-hero.InvincibleStartTime > configs.PropInvincibleTimeMax {
			hero.Invincible = false
			go g.NotifyEntityInfoChange(configs.HeroType, hero.ID, hero, nil)
		}

		// 处理玩家的加速时间
		if hero.SpeedUp && nowTime-hero.SpeedUpStartTime > configs.PropSpeedUpTimeMax {
			originHeroSpeed := configs.HeroSpeedSizeCoeffcient / hero.Size
			if originHeroSpeed < configs.HeroSpeedDownLimit {
				originHeroSpeed = configs.HeroSpeedDownLimit
			}
			hero.Speed = originHeroSpeed
			hero.SpeedUp = false
			go g.NotifyEntityInfoChange(configs.HeroType, hero.ID, hero, nil)
		}

		// 处理玩家的减速时间
		if hero.SpeedDown && nowTime-hero.SpeedDownStartTime > configs.PropSpeedSlowTimeMax {
			originHeroSpeed := configs.HeroSpeedSizeCoeffcient / hero.Size
			if originHeroSpeed > configs.HeroSpeedUpLimit {
				originHeroSpeed = configs.HeroSpeedUpLimit
			}
			hero.Speed = originHeroSpeed
			hero.SpeedDown = false
			go g.NotifyEntityInfoChange(configs.HeroType, hero.ID, hero, nil)
		}

		// 更新玩家位置
		timeElapse := nowTime - hero.UpdateTime
		hero.UpdateTime = nowTime
		distance := float64(timeElapse) * float64(hero.Speed) / 1e9
		x, y := tools.CalXY(distance, hero.HeroDirection.X, hero.HeroDirection.Y)
		hero.HeroPosition.X += x
		hero.HeroPosition.X = tools.GetMax(hero.HeroPosition.X, configs.MapMinX)
		hero.HeroPosition.X = tools.GetMin(hero.HeroPosition.X, configs.MapMaxX)
		hero.HeroPosition.Y += y
		hero.HeroPosition.Y = tools.GetMax(hero.HeroPosition.Y, configs.MapMinY)
		hero.HeroPosition.Y = tools.GetMin(hero.HeroPosition.Y, configs.MapMaxY)
		g.ModifyHero(hero)
	}
	g.onCollision()
}

func (room *GameRoom) onCollision() {
	roomTowers := room.GetTowers()
	// 清空四叉树
	room.quadTree.Clear()
	// 四叉树插入物体
	for _, hero := range room.FetchHeros() {
		// 如果玩家状态为阵亡，则跳过该玩家检测流程
		if hero.Status == int32(pb.HERO_STATUS_DEAD) {
			continue
		}
		//初始化hero加入到四叉树中进行碰撞检测
		room.quadTree.InsertObj(collision.NewRectangleByObj(hero.ID, int32(pb.ENTITY_TYPE_HERO_TYPE), hero.Size, hero.HeroPosition.X, hero.HeroPosition.Y))
	}
	props, _ := room.props.GetProps()
	for _, prop := range props {
		// 如果道具状态为阵亡，则跳过该玩家检测流程
		if prop.Status == int32(pb.ITEM_STATUS_ITEM_DEAD) {
			continue
		}
		// 初始化道具加入到四叉树中进行碰撞检测
		room.quadTree.InsertObj(collision.NewRectangleByObj(prop.Id, prop.PropType, 0, prop.Pos.X, prop.Pos.Y))
	}
	//room.quadTree.Show()
	// 遍历玩家集合，检测碰撞
	heros := room.FetchHeros()
	for heroIndex := 0; heroIndex < len(heros); heroIndex++ {
		hero := heros[heroIndex]
		// 如果玩家状态为阵亡，则跳过该玩家检测流程
		if hero.Status == int32(pb.HERO_STATUS_DEAD) {
			continue
		}
		heroObj := collision.NewRectangleByObj(hero.ID, int32(pb.ENTITY_TYPE_HERO_TYPE), hero.Size, hero.HeroPosition.X, hero.HeroPosition.Y)
		collisionCandidates := room.quadTree.GetObjsInSameDistrict(heroObj)
		for candidateIndex := 0; candidateIndex < len(collisionCandidates); candidateIndex++ {
			// 如果玩家状态为阵亡，则跳过该玩家检测流程
			if hero.Status == int32(pb.HERO_STATUS_DEAD) {
				break
			}
			candidate := collisionCandidates[candidateIndex]
			if collision.CheckCollision(heroObj, candidate) {
				// 检测到发生了碰撞
				// 双方均是英雄，开启碰撞仲裁流程
				if candidate.Type == int32(pb.ENTITY_TYPE_HERO_TYPE) {
					candidateObject, _ := room.Heroes.Load(candidate.ID)
					candidateHero := candidateObject.(*model.Hero)
					if hero.Invincible ||
						candidateHero.Invincible {
						continue
					}
					var loser, winner *model.Hero
					// 仲裁胜负
					if hero.Size > candidate.Size {
						loser = candidateHero
						w, _ := room.Heroes.Load(hero.ID)
						winner = w.(*model.Hero)
						if int32(pb.HERO_STATUS_DEAD) == loser.Status || int32(pb.HERO_STATUS_DEAD) == winner.Status {
							continue
						}

					} else if hero.Size < candidate.Size {
						l, _ := room.Heroes.Load(hero.ID)
						loser = l.(*model.Hero)
						winner = candidateHero
						if int32(pb.HERO_STATUS_DEAD) == loser.Status || int32(pb.HERO_STATUS_DEAD) == winner.Status {
							continue
						}
					}
					if nil == loser || nil == winner {
						continue
					}
					log.Printf("[GameRoom]检测到玩家发生碰撞！胜者信息：%+v，败者信息：%+v\n", winner, loser)
					// 败者退场
					room.Heroes.Delete(loser.ID)
					roomTowers[loser.TowerId].HeroLeave(loser)
					loser.Status = int32(pb.HERO_STATUS_DEAD)
					room.Heroes.Store(loser.ID, loser)
					room.quadTree.DeleteObj(collision.NewRectangleByObj(loser.ID, int32(pb.ENTITY_TYPE_HERO_TYPE), loser.Size, loser.HeroPosition.X, loser.HeroPosition.Y))
					// 胜者增大、变慢、加分
					room.Heroes.Delete(winner.ID)
					winner.Size += candidate.Size
					//eater.Size += configs.HeroSizeGrowthStep
					if winner.Size > configs.HeroSizeUpLimit {
						winner.Size = configs.HeroSizeUpLimit
					}
					winner.Speed = configs.HeroSpeedSizeCoeffcient / winner.Size
					if winner.Speed < configs.HeroSpeedDownLimit {
						winner.Speed = configs.HeroSpeedDownLimit
					}
					winner.Score += configs.HeroEatEnemyBonus

					//更新排行榜
					heroRankInfo := info.NewHeroRankInfo(winner)
					room.heroRankHeap.ChallengeRank(heroRankInfo)
					//发出排行榜变动推送
					rankInfos := room.heroRankHeap.GetSortedHeroRankInfos()
					rankNotify := notify2.NewGameRankListNotify(rankInfos)
					GAME_ROOM_MANAGER.Braodcast(room.ID, rankNotify.ToGMessageBytes())

					room.Heroes.Store(winner.ID, winner)
					room.quadTree.UpdateObj(collision.NewRectangleByObj(winner.ID, int32(pb.ENTITY_TYPE_HERO_TYPE), winner.Size, winner.HeroPosition.X, winner.HeroPosition.Y))
					// 发包
					heroInfo := info.NewHeroInfo(loser)
					notify := notify2.NewEntityInfoChangeNotify(int32(pb.ENTITY_TYPE_HERO_TYPE), loser.ID, heroInfo, nil)
					GAME_ROOM_MANAGER.Braodcast(room.ID, notify.ToGMessageBytes())
					heroInfo = info.NewHeroInfo(winner)
					notify = notify2.NewEntityInfoChangeNotify(int32(pb.ENTITY_TYPE_HERO_TYPE), winner.ID, heroInfo, nil)
					GAME_ROOM_MANAGER.Braodcast(room.ID, notify.ToGMessageBytes())

					//检测是否达到胜利条件
					if winner.Score >= configs.GameWinLiminationScore {
						log.Printf("[GameRoom]英雄达到胜利条件，开始对局结算！hero:%v, room:%v \n", winner, room)
						room.onGameOver()
					}
				}

				// 一方为食物，开启吃道具流程
				if candidate.Type == int32(pb.ENTITY_TYPE_PROP_TYPE_FOOD) ||
					candidate.Type == int32(pb.ENTITY_TYPE_PROP_TYPE_INVINCIBLE) ||
					candidate.Type == int32(pb.ENTITY_TYPE_PROP_TYPE_SPEED_UP) ||
					candidate.Type == int32(pb.ENTITY_TYPE_PROP_TYPE_SIZE_DOWN) ||
					candidate.Type == int32(pb.ENTITY_TYPE_PROP_TYPE_SPEED_DOWN) {
					var prop (*model.Prop)
					var eater (*model.Hero)
					prop, _ = room.props.GetProp(candidate.ID)
					e, _ := room.Heroes.Load(hero.ID)
					eater = e.(*model.Hero)
					if int32(pb.ITEM_STATUS_ITEM_DEAD) == prop.Status || int32(pb.HERO_STATUS_DEAD) == eater.Status {
						continue
					}
					log.Printf("[GameRoom]检测到玩家吃道具！玩家信息：%+v，道具信息：%+v\n", eater, prop)
					// 道具退场
					room.props.RemoveProp(prop.Id)
					prop.Status = int32(pb.ITEM_STATUS_ITEM_DEAD)
					room.props.AddProp(prop)
					roomTowers[prop.TowerId].PropLeave(prop) // 视野管理
					room.quadTree.DeleteObj(collision.NewRectangleByObj(prop.Id, prop.PropType, 0, prop.Pos.X, prop.Pos.Y))
					// 判断吃的道具类型
					switch prop.PropType {
					case int32(pb.ENTITY_TYPE_PROP_TYPE_FOOD):
						// 玩家增大、变慢、加分
						eater.Size += configs.HeroSizeGrowthStep
						if eater.Size > configs.HeroSizeUpLimit {
							eater.Size = configs.HeroSizeUpLimit
						}
						eater.Speed = configs.HeroSpeedSizeCoeffcient / eater.Size
						if eater.Speed < configs.HeroSpeedDownLimit {
							eater.Speed = configs.HeroSpeedDownLimit
						}
						eater.Score += configs.HeroEatItemBonus
						//更新排行榜
						heroRankInfo := info.NewHeroRankInfo(eater)
						room.heroRankHeap.ChallengeRank(heroRankInfo)
						//发出排行榜变动推送
						rankInfos := room.heroRankHeap.GetSortedHeroRankInfos()
						rankNotify := notify2.NewGameRankListNotify(rankInfos)
						GAME_ROOM_MANAGER.Braodcast(room.ID, rankNotify.ToGMessageBytes())
						//检测是否达到胜利条件
						if eater.Score >= configs.GameWinLiminationScore {
							room.onGameOver()
						}

						room.Heroes.Store(eater.ID, eater)
						room.quadTree.UpdateObj(collision.NewRectangleByObj(eater.ID, int32(pb.ENTITY_TYPE_HERO_TYPE), eater.Size, eater.HeroPosition.X, eater.HeroPosition.Y))
						break
					case int32(pb.ENTITY_TYPE_PROP_TYPE_INVINCIBLE):
						eater.InvincibleStartTime = time.Now().UnixNano()
						eater.Invincible = true
						break
					case int32(pb.ENTITY_TYPE_PROP_TYPE_SPEED_UP):
						eater.SpeedUpStartTime = time.Now().UnixNano()
						eater.SpeedUp = true
						eater.SpeedDown = false
						eaterOriginSpeed := configs.HeroSpeedSizeCoeffcient / eater.Size
						eater.Speed = eaterOriginSpeed * 2
						if eater.Speed > configs.HeroSpeedUpLimit {
							eater.Speed = configs.HeroSpeedUpLimit
						}
						break
					case int32(pb.ENTITY_TYPE_PROP_TYPE_SPEED_DOWN):
						eater.SpeedDownStartTime = time.Now().UnixNano()
						eater.SpeedDown = true
						eater.SpeedUp = false
						eaterOriginSpeed := configs.HeroSpeedSizeCoeffcient / eater.Size
						eater.Speed = eaterOriginSpeed / 2
						if eater.Speed < configs.HeroSpeedDownLimit {
							eater.Speed = configs.HeroSpeedSizeCoeffcient
						}
						break
					case int32(pb.ENTITY_TYPE_PROP_TYPE_SIZE_DOWN):
						eater.Size /= 2
						if eater.Size < configs.HeroSizeDownLimit {
							eater.Size = configs.HeroSizeDownLimit
						}
						eater.Speed = configs.HeroSpeedSizeCoeffcient / eater.Size
						if eater.Speed > configs.HeroSpeedUpLimit {
							eater.Speed = configs.HeroSpeedUpLimit
						}
						break
					}
					go room.NotifyEntityInfoChange(prop.PropType, prop.Id, nil, prop)
					//itemInfo := info.NewItemInfo(prop)
					//notify := notify2.NewEntityInfoChangeNotify(prop.PropType, prop.Id, nil, itemInfo)
					//GAME_ROOM_MANAGER.Braodcast(room.ID, notify.ToGMessageBytes())

					go room.NotifyEntityInfoChange(configs.HeroType, eater.ID, eater, nil)
					//heroInfo := info.NewHeroInfo(eater)
					//notify = notify2.NewEntityInfoChangeNotify(int32(pb.ENTITY_TYPE_HERO_TYPE), eater.ID, heroInfo, nil)
					//GAME_ROOM_MANAGER.Braodcast(room.ID, notify.ToGMessageBytes())
				}

			}
		}
	}
}

func (room *GameRoom) onGameOver() {
	//广播游戏结算推送
	log.Printf("[GameRoom]游戏对局结束，开始广播至玩家并回收资源！")
	heroRankInfos := room.heroRankHeap.GetSortedHeroRankInfos()
	notify := notify2.NewGameFinishNotify(heroRankInfos, time.Now().Unix())
	GAME_ROOM_MANAGER.Braodcast(room.ID, notify.ToGMessageBytes())
	//回收游戏对局资源
	//todo
	atomic.AddInt32(&room.gameOver, 1) // 发通知，告诉 goroutine 游戏结束
	GAME_ROOM_MANAGER.DeleteGameRoom(room.ID)
	//runtime.GC()
}
