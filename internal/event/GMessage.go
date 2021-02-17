package event

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
)

type GMessage struct {
	framework.BaseEvent
	MsgType     int32
	GameMsgCode int32
	SessionId   int32
	SeqId       int32
	Data        e.Event
}

func (this *GMessage) ToMessage() interface{} {
	var pbMsgNotify *pb.Notify
	var pbMsgRequest *pb.Request
	var pbMsgResponse *pb.Response

	switch this.MsgType {
	case tools.MsgTypeNotify:
		switch this.Data.(type) { //这里的处理函数可以进行封装
		case *notify.EntityInfoChangeNotify:
			pbMsgNotify = &pb.Notify{
				EntityInfoChangeNotify: this.Data.ToMessage().(*pb.EntityInfoChangeNotify),
			}

		default:
			panic("no match type")
		}
		break
	case tools.MsgTypeRequest:
		switch this.Data.(type) {
		case *request.EntityInfoChangeRequest:
			pbMsgRequest = &pb.Request{
				EntityChangeRequest: this.Data.ToMessage().(*pb.EntityInfoChangeRequest),
			}
		}
		break
	case tools.MsgTypeResponse:
		break
	default:
		panic("msg type is incorrect")
	}
	pbMsg := &pb.GMessage{
		MsgType:   pb.MSG_TYPE(this.MsgType),
		MsgCode:   pb.GAME_MSG_CODE(this.GameMsgCode),
		SessionId: this.SessionId,
		SeqId:     this.SeqId,
		Notify:    pbMsgNotify,
		Request:   pbMsgRequest,
		Response:  pbMsgResponse,
	}

	return pbMsg
}

func (this *GMessage) CopyFromMessage(obj interface{}) e.Event {
	pbMsg := obj.(*pb.GMessage)
	msg := &GMessage{
		MsgType: int32(pbMsg.MsgType),
		SeqId:   pbMsg.SeqId,
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

func (this *GMessage) FromMessage(obj interface{}) {
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
