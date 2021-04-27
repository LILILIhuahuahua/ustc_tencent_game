package game

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/kcpnet"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	notify2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/scheduler"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/golang/protobuf/proto"
	"log"
)

type GameRoomManager struct {
	addr    string
	server  *kcpnet.KcpServer
	roomMap map[int64]*GameRoom
	timer   *scheduler.Timer
	waitSessionMap map[int32]*framework.BaseSession //会话等待集合，匹配到一定玩家会话才进行派遣
}

var GAME_ROOM_MANAGER *GameRoomManager

func init() {
	GAME_ROOM_MANAGER = NewGameRoomManager(configs.ServerAddr)
	GAME_ROOM_MANAGER.timer = scheduler.NewTimer(configs.GlobalInfoNotifyInterval, GlobalInfoNotify)
}

func NewGameRoomManager(address string) *GameRoomManager {
	//s, _ := kcpnet.NewKcpServer(address)
	return &GameRoomManager{
		addr:    address,
		server:  nil,
		roomMap: make(map[int64]*GameRoom),
		waitSessionMap: make(map[int32]*framework.BaseSession),
	}
}

func (manager *GameRoomManager) Serv() error {
	//for {
	//	conn, err := manager.server.Listen.AcceptKCP()
	//	if err != nil {
	//		return err
	//	}
	//	conn.SetWindowSize(4800, 4800)
	//	session := framework.NewBaseSession(conn)
	//	manager.dispatchSessionToGameRome(session)
	//}
	return nil
}

func (manager *GameRoomManager) dispatchSessionToGameRome(session *framework.BaseSession) {
	//1.等待一批足够数量的会话
	//2.等待时间超过某个值
	//3.寻找一个可加入的房间
}

func (manager *GameRoomManager) fetchAccessableGameRome() *GameRoom {
	return nil
}

// Cron initialize timer for GameRoomManager. It will broadcast props/food info of each room to its clients.
func GlobalInfoNotify() {
	// set cron function for room
	//log.Printf("Len of roomMap: %d", len(GAME_ROOM_MANAGER.roomMap))
	if GAME_ROOM_MANAGER != nil && len(GAME_ROOM_MANAGER.roomMap) != 0 {
		for _, room := range GAME_ROOM_MANAGER.roomMap {
			var pbMsg *pb.GMessage

			roomHeros := room.GetHeros()
			// 该房间内所有的小球
			var heroNeedToNotify []*model.Hero
			roomHeros.Range(func(k, v interface{}) bool {
				heroNeedToNotify = append(heroNeedToNotify, v.(*model.Hero))
				return true
			})
			for _, hero := range heroNeedToNotify {
				var heroInfos []info.HeroInfo
				players, props := room.GetItemsNearby(hero)
				var items []info.ItemInfo
				for _, v := range props {
					itemInfo := info.ItemInfo{
						ID:     v.Id,
						Type:   v.PropType,
						Status: v.Status,
						ItemPosition: &info.CoordinateXYInfo{
							CoordinateX: v.Pos.X,
							CoordinateY: v.Pos.Y,
						},
					}
					items = append(items, itemInfo)
				}
				for _, player := range players {
					//heroInfo := player.ToEvent()
					heroInfo := *info.NewHeroInfo(player)
					heroInfos = append(heroInfos, heroInfo)
				}
				heroNum := int32(len(players))
				notify := notify2.GameGlobalInfoNotify{
					HeroNumber: heroNum,
					//Time:       0,
					HeroInfos: heroInfos,
					ItemInfos: items,
					//MapInfo:    info.MapInfo{},
				}
				msg := event2.GMessage{
					MsgType:     configs.MsgTypeNotify,
					GameMsgCode: configs.GameGlobalInfoNotify,
					//SessionId:   this.room.,
					Data: &notify,
				}

				pbMsg = msg.ToMessage().(*pb.GMessage)
				data, err := proto.Marshal(pbMsg)
				if err != nil {
					log.Printf("fail to marshal: %s", err.Error())
				}
				GAME_ROOM_MANAGER.Unicast(room.GetRoomID(), hero.Session.Id, data)
			}

		}
	}
}

func (m *GameRoomManager) FetchGameRoom(id int64) *GameRoom {
	return m.roomMap[id]
}

func (m *GameRoomManager) RegisterGameRoom(room *GameRoom) {
	if nil == m.roomMap[room.GetRoomID()] {
		m.roomMap[room.GetRoomID()] = room
	}
}

func (m *GameRoomManager) Unicast(roomId int64, sessionId int32, buff []byte) {
	defer func() {
		e := recover()
		if e != nil {
			//fmt.Println("在unicast的时候出错了，错误为：", e)
		}
	}()

	r := m.FetchGameRoom(roomId)
	s := r.FetchConnector(sessionId)
	if s == nil {
		panic("没有该玩家")
	}
	err := s.SendMessage(buff)
	if err != nil {
		panic(err)
	}
	//m.FetchGameRoom(roomId).FetchConnector(sessionId).SendMessage(buff)
}

func (m *GameRoomManager) Braodcast(roomId int64, buff []byte) {
	r := m.FetchGameRoom(roomId)
	r.BroadcastAll(buff)
	//m.FetchGameRoom(roomId).FetchConnector(sessionId).SendMessage(buff)
}

func (m *GameRoomManager) MutiplecastToNearBy(roomId int64, buf []byte, hero *model.Hero) {
	r := m.FetchGameRoom(roomId)
	var sessionToSend []*framework.BaseSession
	heroToSend, _ := r.GetItemsNearby(hero)
	for _, hero := range heroToSend {
		if hero.Session != nil && hero.Session.Status == configs.SessionStatusCreated {
			sessionToSend = append(sessionToSend, hero.Session)
		} else {
			fmt.Println("hero的session为null")
		}
	}
	r.Multiplecast(buf, sessionToSend)
}

func (m *GameRoomManager) DeleteUnavailableSession() {
	for _, room := range m.roomMap {
		err := room.DeleteUnavailableSession()
		if err != nil {
			fmt.Println("清理不可用session的时候发生了error: ", err.Error())
		}
	}
}

func (m *GameRoomManager) DeleteDeprecatedHero() {
	for _, room := range m.roomMap {
		err := room.DeleteOfflinePlayer()
		if err != nil {
			fmt.Println("清理废弃玩家的小球时发生了error: ", err.Error())
		}
	}
}

func (m *GameRoomManager) UpdateHeroPositionAndStatus() {
	for _, room := range m.roomMap {
		room.UpdateHeroPosAndStatus()
	}
}
