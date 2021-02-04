package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type EnterGameResponse struct {
	framework.BaseEvent
	Result bool
	HeroId int32
}

func (e *EnterGameResponse) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EnterGameResponse)
	e.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE))
	e.Result = pbMsg.GetChangeResult()
	e.HeroId = pbMsg.GetHeroId()
}

func (e *EnterGameResponse) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Response).EnterGameResponse
	resp := &EnterGameResponse{
		Result: pbMsg.GetChangeResult(),
		HeroId: pbMsg.GetHeroId(),
	}
	resp.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE))
	return resp
}

func (e *EnterGameResponse) ToMessage() interface{} {
	return pb.EnterGameResponse{
		ChangeResult: e.Result,
		HeroId:       e.HeroId,
	}
}
