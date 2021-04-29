package request

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
)

type EntityInfoChangeRequest struct {
	framework.BaseEvent //基础消息类作为父类
	EventType           int32
	HeroId              int32
	LinkedId            int32
	LinkedType          string
	HeroMsg             info.HeroInfo
}

func (e *EntityInfoChangeRequest) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.Request).EntityChangeRequest
	e.Code = int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)
	e.EventType = int32(pbMsg.EventType)
	e.HeroId = pbMsg.HeroId
	e.LinkedId = pbMsg.LinkedId
	e.LinkedType = pbMsg.LinkedType.String()
	info := info.HeroInfo{}
	info.FromMessage(pbMsg.GetHeroMsg())
	e.HeroMsg = info
}

func (e *EntityInfoChangeRequest) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Request).EntityChangeRequest
	info := info.HeroInfo{}
	info.FromMessage(pbMsg.GetHeroMsg())
	req := EntityInfoChangeRequest{
		EventType:  int32(pbMsg.EventType),
		HeroId:     pbMsg.HeroId,
		LinkedId:   pbMsg.LinkedId,
		LinkedType: pbMsg.LinkedType.String(),
		HeroMsg:    info,
	}
	req.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	return &req
}

func (e *EntityInfoChangeRequest) ToMessage() interface{} {
	//todo:改这里写死的类型
	return &pb.EntityInfoChangeRequest{
		EventType:  pb.EVENT_TYPE_HERO_MOVE,
		HeroId:     e.HeroId,
		LinkedId:   e.LinkedId,
		LinkedType: pb.ENTITY_TYPE_HERO_TYPE,
		HeroMsg:    e.HeroMsg.ToMessage().(*pb.HeroMsg),
	}
}
