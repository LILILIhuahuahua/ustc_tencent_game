package request

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type EntityInfoChangeRequest struct {
	framework.BaseEvent	//基础消息类作为父类
	EventType  string
	HeroId     int32
	LinkedId   int32
	LinkedType string
	//heroMsg
}

func (this *EntityInfoChangeRequest)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.Request).EntityChangeRequest
	this.Code = int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)
	this.EventType = pbMsg.EventType.String()
	this.HeroId = pbMsg.HeroId
	this.LinkedId = pbMsg.LinkedId
	this.LinkedType = pbMsg.LinkedType.String()
}

func (this *EntityInfoChangeRequest)CopyFromMessage(obj interface{}) event.Event{
	pbMsg := obj.(*pb.Request).EntityChangeRequest
	req :=  EntityInfoChangeRequest{
		//Code:  int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST),
		EventType:  pbMsg.EventType.String(),
		HeroId:     pbMsg.HeroId,
		LinkedId:   pbMsg.LinkedId,
		LinkedType: pbMsg.LinkedType.String(),
	}
	req.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	return &req
}

func (this *EntityInfoChangeRequest)ToMessage() interface{} {
	return &pb.EntityInfoChangeRequest{
		EventType: pb.EVENT_TYPE_HERO_MOVE,
		HeroId: this.HeroId,
		LinkedId: this.LinkedId,
		LinkedType: pb.ENTITY_TYPE_HERO_TYPE,
	}
}
