package game

import (
	"fmt"
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event/info"
	notify2 "github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event/request"
	response2 "github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/golang/protobuf/proto"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type GameEventHandler struct{}

var GAME_EVENT_HANDLER = &GameEventHandler{}

func (g GameEventHandler) OnEvent(e event.Event) {
	if nil == e {
		return
	}
	// 二级解码
	msg := e.(*event2.GMessage)
	data := msg.Data
	switch data.GetCode() {

	case int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST):
		g.onEntityInfoChange(data.(*request.EntityInfoChangeRequest))

	case int32(pb.GAME_MSG_CODE_HEART_BEAT_REQUEST):
		g.onHeartBeat(data.(*request.HeartBeatRequest))

	case int32(pb.GAME_MSG_CODE_HERO_QUIT_REQUEST):
		g.onHeroQuit(data.(*request.HeroQuitRequest))

	default:
		return
	}
}

func (g GameEventHandler) OnEventToSession(e event.Event, s event.Session) {

}

func (g GameEventHandler) onHeartBeat(req *request.HeartBeatRequest) {
	sendTime := req.SendTime
	heartBeatRsp := response2.NewHeartBeatResponse(sendTime)
	outResponse := heartBeatRsp.ToGMessageBytes(req.SeqId)
	GAME_ROOM_MANAGER.Unicast(req.GetRoomId(), req.SessionId, outResponse)
}

func (g GameEventHandler) onHeroQuit(req *request.HeroQuitRequest) {
	heroID := req.HeroId
	sessionId := req.SessionId
	//1.更改玩家状态为dead
	roomID := req.GetRoomId()
	room := GAME_ROOM_MANAGER.FetchGameRoom(roomID)
	h, _ := room.Heroes.Load(heroID)
	if nil == h {
		fmt.Println("[HeroQuitErr]无法找到对应英雄！")
	}
	hero := h.(*model.Hero)
	room.Heroes.Delete(heroID)
	hero.Status = int32(pb.HERO_STATUS_DEAD)
	room.Heroes.Store(heroID, hero)
	log.Printf("[GameEventHandler]玩家离开游戏房间！hero：%v, room:%v\n", hero, room)
	atomic.AddInt32(&room.AliveHeroNum, -1)
	//2.广播给其他玩家
	heroInfo := info.NewHeroInfo(hero)
	notify := notify2.NewEntityInfoChangeNotify(int32(pb.ENTITY_TYPE_HERO_TYPE), hero.ID, heroInfo, nil)
	out := notify.ToGMessageBytes()
	GAME_ROOM_MANAGER.Braodcast(room.ID, out)
	//3.单播返回离开结果
	resp := response2.NewHeroQuitResponse(true)
	out = resp.ToGMessageBytes(req.SeqId)
	GAME_ROOM_MANAGER.Unicast(roomID, sessionId, out)
}

func (g GameEventHandler) onEntityInfoChange(req *request.EntityInfoChangeRequest) {
	r := GAME_ROOM_MANAGER.FetchGameRoom(req.RoomId)
	var pbNotifyMsg, pbResponseMsg *pb.GMessage

	if req.EventType == int32(pb.EVENT_TYPE_HERO_MOVE) {
		if r.GetHero(req.HeroMsg.ID) == nil {
			fmt.Println("hero为nil，不ok")
			return
		}
		if req.HeroMsg.Speed == float32(0) {
			req.HeroMsg.Speed = float32(100)
		}
		//fmt.Printf("我收到的X为%f, Y为%f", req.HeroMsg.HeroDirection.CoordinateX, req.HeroMsg.HeroDirection.CoordinateY)
		hero := &model.Hero{
			ID:            req.HeroMsg.ID,
			Status:        req.HeroMsg.Status,
			Size:          req.HeroMsg.Size,
			Speed:         req.HeroMsg.Speed,
			UpdateTime:    time.Now().UnixNano(),
			HeroDirection: model.Coordinate{X: req.HeroMsg.HeroDirection.CoordinateX, Y: req.HeroMsg.HeroDirection.CoordinateY},
			HeroPosition:  model.Coordinate{X: req.HeroMsg.HeroPosition.CoordinateX, Y: req.HeroMsg.HeroPosition.CoordinateY},
		}
		//heros := room.GetHeros()
		//heroObj, ok := heros.Load(req.HeroMsg.ID)
		//if !ok {
		//	panic("hero not exists") //之后改成response false
		//}
		//nowHero := heroObj.(*model.Hero)
		//if !tools.JudgePosition(
		//	req.HeroMsg.HeroPosition.CoordinateX,
		//	req.HeroMsg.HeroPosition.CoordinateY,
		//	nowHero.HeroDirection.X,
		//	nowHero.HeroDirection.Y,
		//	nowHero.Speed) {
		//	panic("hero position is not correct")
		//}

		//todo:封装解包类方法
		var lock = sync.Mutex{}
		lock.Lock()
		r.ModifyHero(hero)
		hero = r.GetHero(hero.ID)
		lock.Unlock()

		notify := &notify2.EntityInfoChangeNotify{
			EntityType: configs.HeroType,
			EntityId:   hero.ID,
			//HeroMsg:    hero.ToEvent(),
			HeroMsg: info.NewHeroInfo(hero),
			//ItemMsg: nil,
		}

		response := &response2.EntityInfoChangeResponse{
			ChangeResult: true,
		}

		notifyMsg := event2.GMessage{
			MsgType:     configs.MsgTypeNotify,
			GameMsgCode: configs.EntityInfoNotify, //这里命名之后要修改
			SessionId:   req.SessionId,
			Data:        notify,
		}

		responseMsg := event2.GMessage{
			MsgType:     configs.MsgTypeResponse,
			GameMsgCode: configs.EntityInfoChangeResponse,
			SessionId:   req.SessionId,
			SeqId:       req.SeqId,
			Data:        response,
		}
		pbNotifyMsg = notifyMsg.ToMessage().(*pb.GMessage)
		//fmt.Printf("发送的消息为%v \n", pbNotifyMsg)
		pbResponseMsg = responseMsg.ToMessage().(*pb.GMessage)
		outNotify, err := proto.Marshal(pbNotifyMsg)
		outResponse, err := proto.Marshal(pbResponseMsg)
		if nil == err {
		}
		//GAME_ROOM_MANAGER.Braodcast(req.GetRoomId(), outNotify)
		GAME_ROOM_MANAGER.MutiplecastToNearBy(r.ID, outNotify, hero) // 只通知视野范围内的玩家,而非广播给所有的玩家
		GAME_ROOM_MANAGER.Unicast(req.GetRoomId(), req.SessionId, outResponse)
	}
}
