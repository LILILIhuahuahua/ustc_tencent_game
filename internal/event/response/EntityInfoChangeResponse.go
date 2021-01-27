package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)
type EntityInfoChangeResponse struct {
	Code         int32
	ChangeResult bool
}

func (this EntityInfoChangeResponse) GetCode() int32 {
	return this.Code
}

func (this EntityInfoChangeResponse) SetCode(eventCode int32) {
	this.Code = eventCode
}

func (this EntityInfoChangeResponse)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EntityInfoChangeResponse)
	this.Code = int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)
	this.ChangeResult = pbMsg.GetChangeResult()
}

func (this EntityInfoChangeResponse)CopyFromMessage(obj interface{}) event.Event{
	pbMsg := obj.(*pb.EntityInfoChangeResponse)
	return EntityInfoChangeResponse{
		Code:         int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST),
		ChangeResult: pbMsg.GetChangeResult(),
	}
}

func (this EntityInfoChangeResponse)ToMessage() interface{} {
	return pb.EntityInfoChangeResponse{
		ChangeResult: this.ChangeResult,
	}
}

