package notify

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
)

type HeroViewNotify struct {
	framework.BaseEvent
	HeroId   int32
	ViewType int32
	HeroMsg  info.HeroInfo
}

func (notify *HeroViewNotify) FromMessage(obj interface{}) {

}

func (notify *HeroViewNotify) CopyFromMessage(obj interface{}) e.Event {
	return &HeroViewNotify{}
}

func (notify *HeroViewNotify) ToMessage() interface{} {
	pbMsg := &pb.HeroViewNotify{
		HeroId:   notify.HeroId,
		ViewType: pb.VIEW_TYPE(notify.ViewType),
		HeroMsg:  notify.HeroMsg.ToMessage().(*pb.HeroMsg),
	}
	return pbMsg
}
