package handler

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/golang/protobuf/proto"
)

type GameEventHandler struct {
}

var GAME_EVENT_HANDLER = &GameEventHandler{}

func (this GameEventHandler)OnEvent(e event.Event) {
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

func (this GameEventHandler)OnEventToSession(e event.Event,s event.Session) {

}

func (this GameEventHandler)onEntityInfoChange(req *request.EntityInfoChangeRequest) {
	id := req.HeroId
	println(id)
	//回包
	data := pb.EntityInfoChangeResponse{
		ChangeResult: true,
	}
	resp:=pb.Response{
		EntityChangeResponse: &data,
	}
	msg := pb.GMessage{
		MsgType: pb.MSG_TYPE_RESPONSE,
		MsgCode: pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE,
		Response: &resp,
	}
	out, err := proto.Marshal(&msg)
	if nil==err {

	}
	framework.GAME_ROOM_MANAGER.Unicast(req.GetRoomId(), req.GetSessionId(), out)
}
