package event

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)
type GMessage struct {
	framework.BaseEvent
	MsgType int32
	//GameMsgCode int32
	//SessionId int64
	SeqId int32
	Data e.Event
}

func (this *GMessage) ToMessage() interface{} {
	//todo:
	panic("implement me")
}

func (this *GMessage) CopyFromMessage(obj interface{}) e.Event {
	pbMsg := obj.(*pb.GMessage)
	msg := &GMessage{
		MsgType: int32(pbMsg.MsgType),
		SeqId: pbMsg.SeqId,
	}
	msg.SetCode(int32(pbMsg.MsgCode))
	msg.SetSessionId(pbMsg.SessionId)
	event := e.Manager.FetchEvent(msg.GetCode())
	if pb.MSG_TYPE_NOTIFY == pbMsg.MsgType {
		msg.Data = event.CopyFromMessage(pbMsg.Notify)
	}
	if pb.MSG_TYPE_REQUEST == pbMsg.MsgType {
		msg.Data = event.CopyFromMessage(pbMsg.Request)
	}
	if pb.MSG_TYPE_RESPONSE == pbMsg.MsgType {
		msg.Data = event.CopyFromMessage(pbMsg.Response)
	}
	//传递会话id至二层协议中
	msg.Data.SetSessionId(pbMsg.SessionId)
	msg.Data.SetRoomId(this.RoomId)
	return msg
}

//todo：值类型的receiver只能使用值传递
func (this *GMessage)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.GMessage)
	this.MsgType = int32(pbMsg.MsgType)
	this.Code = int32(pbMsg.MsgCode)
	this.SessionId = pbMsg.SessionId
	this.SeqId = pbMsg.SeqId

	event := e.Manager.FetchEvent(this.GetCode())
	if pb.MSG_TYPE_NOTIFY == pbMsg.MsgType {
		this.Data = event.CopyFromMessage(pbMsg.Notify)
	}
	if pb.MSG_TYPE_REQUEST == pbMsg.MsgType {
		this.Data = event.CopyFromMessage(pbMsg.Request)
	}
	if pb.MSG_TYPE_RESPONSE == pbMsg.MsgType {
		this.Data = event.CopyFromMessage(pbMsg.Response)
	}
	this.Data.SetSessionId(pbMsg.SessionId)
	this.Data.SetRoomId(this.RoomId)
}