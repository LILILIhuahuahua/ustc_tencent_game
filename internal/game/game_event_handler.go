package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	notify2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	response2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/golang/protobuf/proto"
	"sync"
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

	default:
		return
	}
}

func (g GameEventHandler) OnEventToSession(e event.Event, s event.Session) {

}

func (g GameEventHandler) onEntityInfoChange(req *request.EntityInfoChangeRequest) {
	//heroId := req.HeroId
	//room := GAME_ROOM_MANAGER.FetchGameRoom(req.RoomId)
	r := GAME_ROOM_MANAGER.FetchGameRoom(req.RoomId)
	var pbNotifyMsg, pbResponseMsg *pb.GMessage

	hero := &model.Hero{}
	hero.ID = req.HeroMsg.ID
	hero.Speed = req.HeroMsg.Speed
	hero.Size = req.HeroMsg.Size
	hero.Status = req.HeroMsg.Status
	hero.HeroPosition = model.Coordinate{}
	hero.HeroPosition.X = req.HeroMsg.HeroPosition.CoordinateX
	hero.HeroPosition.Y = req.HeroMsg.HeroPosition.CoordinateY
	hero.HeroDirection = model.Coordinate{}
	hero.HeroDirection.X = req.HeroMsg.HeroDirection.CoordinateX
	hero.HeroDirection.Y = req.HeroMsg.HeroDirection.CoordinateY

	if req.EventType == int32(pb.EVENT_TYPE_HERO_MOVE) {
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
		lock.Unlock()

		notify := &notify2.EntityInfoChangeNotify{
			EntityType: configs.HeroType,
			EntityId:   hero.ID,
			HeroMsg:    hero.ToEvent(),
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
			MsgType:	 configs.MsgTypeResponse,
			GameMsgCode: configs.EntityInfoChangeResponse,
			SessionId: 	 req.SessionId,
			Data:		 response,
		}
		pbNotifyMsg = notifyMsg.ToMessage().(*pb.GMessage)
		pbResponseMsg = responseMsg.ToMessage().(*pb.GMessage)
		outNotify, err := proto.Marshal(pbNotifyMsg)
		outResponse, err := proto.Marshal(pbResponseMsg)
		if nil == err {
		}
		GAME_ROOM_MANAGER.Braodcast(req.GetRoomId(), outNotify)
		GAME_ROOM_MANAGER.Unicast(req.GetRoomId(), req.SessionId ,outResponse)
	} else if req.EventType == int32(pb.EVENT_TYPE_HERO_COLLISION) {
		collisionRes := true
		//todo:碰撞检测
		response := &response2.EntityInfoChangeResponse{
			ChangeResult: collisionRes,
		}
		responseMsg := event2.GMessage{
			MsgType:	 configs.MsgTypeResponse,
			GameMsgCode: configs.EntityInfoChangeResponse,
			SessionId: 	 req.SessionId,
			Data:		 response,
		}
		pbResponseMsg = responseMsg.ToMessage().(*pb.GMessage)
		outResponse, _ := proto.Marshal(pbResponseMsg)
		GAME_ROOM_MANAGER.Unicast(req.GetRoomId(), req.SessionId ,outResponse)
	}
}
