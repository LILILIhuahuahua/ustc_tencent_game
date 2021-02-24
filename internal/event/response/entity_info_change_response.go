package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type EntityInfoChangeResponse struct {
	framework.BaseEvent //基础消息类作为父类
	ChangeResult        bool
}

func (this *EntityInfoChangeResponse) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EntityInfoChangeResponse)
	this.Code = int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)
	this.ChangeResult = pbMsg.GetChangeResult()
}

func (this *EntityInfoChangeResponse) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.EntityInfoChangeResponse)
	resp := &EntityInfoChangeResponse{
		ChangeResult: pbMsg.GetChangeResult(),
	}
	resp.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	return resp
}

func (this *EntityInfoChangeResponse) ToMessage() interface{} {
	return &pb.EntityInfoChangeResponse{
		ChangeResult: this.ChangeResult,
	}
}
