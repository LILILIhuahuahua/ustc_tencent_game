package handler

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
)

type GameEventHandler struct {
}

func (this GameEventHandler)OnEvent(e event.Event) {
	if nil == e {
		return
	}
	switch e.GetCode() {
		case int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST):
			this.onEntityInfoChange(e.(request.EntityInfoChangeRequest))
		default:
			return
	}
}

func (this GameEventHandler)OnEventToSession(e event.Event,s event.Session) {

}

func (this GameEventHandler)onEntityInfoChange(req request.EntityInfoChangeRequest) {
	id := req.HeroId
	println(id)
}
