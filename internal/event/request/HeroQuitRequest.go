package request

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type HeroQuitRequest struct {
	framework.BaseEvent
	HeroId   int32
	QuitTime int64
}

func (e *HeroQuitRequest) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.HeroQuitRequest)
	e.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_REQUEST))
	e.HeroId = pbMsg.GetHeroId()
	e.QuitTime = pbMsg.GetTime()
}

func (e *HeroQuitRequest) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Request).HeroQuitRequest
	req := &HeroQuitRequest{
		HeroId:   pbMsg.GetHeroId(),
		QuitTime: pbMsg.GetTime(),
	}
	req.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_REQUEST))
	return req
}

func (e *HeroQuitRequest) ToMessage() interface{} {
	return pb.HeroQuitRequest{
		HeroId: e.HeroId,
		Time:   e.QuitTime,
	}
}
