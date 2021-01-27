package event

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)
type GMessage struct {
	MsgType int32
	GameMsgCode int32
	SessionId int64
	SeqId int32
	Data e.Event
}

func (this GMessage) GetCode() int32 {
	return this.GameMsgCode
}

func (this GMessage) SetCode(eventCode int32) {
	panic("implement me")
}

func (this GMessage) ToMessage() interface{} {
	panic("implement me")
}

func (this GMessage) CopyFromMessage(obj interface{}) e.Event {
	pbMsg := obj.(*pb.GMessage)
	msg := GMessage{
		MsgType: int32(pbMsg.MsgType),
		GameMsgCode: int32(pbMsg.MsgCode),
		SessionId: pbMsg.SessionId,
		SeqId: pbMsg.SeqId,
	}
	event := e.Manager.FetchEvent(msg.GameMsgCode)
	if pb.MSG_TYPE_NOTIFY == pbMsg.MsgType {
		msg.Data = event.CopyFromMessage(pbMsg.Notify)
	}
	if pb.MSG_TYPE_REQUEST == pbMsg.MsgType {
		msg.Data = event.CopyFromMessage(pbMsg.Request)
	}
	if pb.MSG_TYPE_RESPONSE == pbMsg.MsgType {
		msg.Data = event.CopyFromMessage(pbMsg.Response)
	}
	return msg
}

//todo：值类型的receiver只能使用值传递
func (this GMessage)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.GMessage)
	this.MsgType = int32(pbMsg.MsgType)
	this.GameMsgCode = int32(pbMsg.MsgCode)
	this.SessionId = pbMsg.SessionId
	this.SeqId = pbMsg.SeqId

	event := e.Manager.FetchEvent(this.GameMsgCode)
	if pb.MSG_TYPE_NOTIFY == pbMsg.MsgType {
		this.Data = event.CopyFromMessage(pbMsg.Notify)
	}
	if pb.MSG_TYPE_REQUEST == pbMsg.MsgType {
		this.Data = event.CopyFromMessage(pbMsg.Request)
	}
	if pb.MSG_TYPE_RESPONSE == pbMsg.MsgType {
		this.Data = event.CopyFromMessage(pbMsg.Response)
	}

}