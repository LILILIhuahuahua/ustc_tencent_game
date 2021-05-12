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

func (e *EntityInfoChangeResponse) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EntityInfoChangeResponse)
	e.Code = int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)
	e.ChangeResult = pbMsg.GetChangeResult()
}

func (e *EntityInfoChangeResponse) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.EntityInfoChangeResponse)
	resp := &EntityInfoChangeResponse{
		ChangeResult: pbMsg.GetChangeResult(),
	}
	resp.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	return resp
}

func (e *EntityInfoChangeResponse) ToMessage() interface{} {
	return &pb.EntityInfoChangeResponse{
		ChangeResult: e.ChangeResult,
	}
}
