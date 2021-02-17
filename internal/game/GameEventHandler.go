package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	notify2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/model"
	"github.com/golang/protobuf/proto"
	"sync"
)

type GameEventHandler struct {}

var GAME_EVENT_HANDLER = &GameEventHandler{}

func (this GameEventHandler) OnEvent(e event.Event) {
	if nil == e {
		return
	}
	// 二级解码
	msg := e.(*event2.GMessage)
	data := msg.Data
	switch data.GetCode() {

	case int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST):
		this.onEntityInfoChange(data.(*request.EntityInfoChangeRequest))

	default:
		return
	}
}

func (this GameEventHandler) OnEventToSession(e event.Event, s event.Session) {

}

func (this GameEventHandler) onEntityInfoChange(req *request.EntityInfoChangeRequest) {
	//heroId := req.HeroId
	g := GAME_ROOM_MANAGER.FetchGameRoom(req.RoomId)
	var pbMsg *pb.GMessage
	if req.EventType == int32(pb.EVENT_TYPE_HERO_MOVE) {
		//todo:封装解包类方法
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
		var lock = sync.Mutex{}
		lock.Lock()
		g.ModifyHero(hero)
		lock.Unlock()

		notify := &notify2.EntityInfoChangeNotify{
			EntityType: configs.HeroType,
			EntityId:   hero.ID,
			HeroMsg:    hero.ToEvent(),
			//ItemMsg: nil,
		}

		msg := event2.GMessage{
			MsgType:     configs.MsgTypeNotify,
			GameMsgCode: configs.EntityInfoNotify, //这里命名之后要修改
			SessionId:   req.SessionId,
			Data:        notify,
		}
		pbMsg = msg.ToMessage().(*pb.GMessage)
	}

	//回包
	//heros := g.FetchHeros()
	//data := pb.GameGlobalInfoNotify{
	//	HeroNumber: int32(len(heros)),
	//}
	//for _, h := range heros {
	//	hMsg := &pb.HeroMsg{}
	//	hMsg.HeroId = h.ID
	//	hMsg.HeroSize = h.Size
	//	hMsg.HeroSpeed = h.Speed
	//	hMsg.HeroPosition = &pb.CoordinateXY{}
	//	hMsg.HeroPosition.CoordinateX = h.HeroPosition.X
	//	hMsg.HeroPosition.CoordinateY = h.HeroPosition.Y
	//	hMsg.HeroDirection = &pb.CoordinateXY{}
	//	hMsg.HeroDirection.CoordinateX = h.HeroDirection.X
	//	hMsg.HeroDirection.CoordinateY = h.HeroDirection.Y
	//	data.HeroMsg = append(data.HeroMsg, hMsg)
	//}
	//


	//data := pb.EntityInfoChangeResponse{
	//	ChangeResult: true,
	//}
	//resp:=pb.Response{
	//	EntityChangeResponse: &data,
	//}
	//msg := pb.GMessage{
	//	MsgType: pb.MSG_TYPE_RESPONSE,
	//	MsgCode: pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE,
	//	Response: &resp,
	//}
	out, err := proto.Marshal(pbMsg)
	if nil == err {

	}
	GAME_ROOM_MANAGER.Braodcast(req.GetRoomId(), out)
}
