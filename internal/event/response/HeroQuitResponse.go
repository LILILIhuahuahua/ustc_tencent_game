package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type HeroQuitResponse struct {
	framework.BaseEvent
	QuitResult bool
}

func (e *HeroQuitResponse) ToMessage() interface{} {
	return pb.HeroQuitResponse{
		QuitResult: e.QuitResult,
	}
}


func (e *HeroQuitResponse) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.HeroQuitResponse)
	e.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_RESPONSE))
	e.QuitResult = pbMsg.GetQuitResult()
}

func (e *HeroQuitResponse) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Response).HeroQuitResponse
	resp := &HeroQuitResponse{
		QuitResult: pbMsg.GetQuitResult(),
	}
	resp.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_RESPONSE))
	return resp
}