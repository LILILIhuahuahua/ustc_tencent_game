package event

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	response2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
)

type GMessage struct {
	framework.BaseEvent
	MsgType     int32
	GameMsgCode int32
	SessionId   int32
	SeqId       int32
	SendTime    int64
	Data        e.Event
}

func (g *GMessage) ToMessage() interface{} {
	var pbMsgNotify *pb.Notify
	var pbMsgRequest *pb.Request
	var pbMsgResponse *pb.Response

	switch g.MsgType {
	case configs.MsgTypeNotify:
		switch g.Data.(type) { //这里的处理函数可以进行封装
		case *notify.EntityInfoChangeNotify:
			pbMsgNotify = &pb.Notify{
				EntityInfoChangeNotify: g.Data.ToMessage().(*pb.EntityInfoChangeNotify),
			}
		case *notify.GameGlobalInfoNotify:
			pbMsgNotify = &pb.Notify{
				GameGlobalInfoNotify: g.Data.ToMessage().(*pb.GameGlobalInfoNotify),
			}
		case *notify.EnterGameNotify:
			pbMsgNotify = &pb.Notify{
				EnterGameNotify: g.Data.ToMessage().(*pb.EnterGameNotify),
			}
		case *notify.HeroViewNotify:
			pbMsgNotify = &pb.Notify{
				HeroViewNotify: g.Data.ToMessage().(*pb.HeroViewNotify),
			}
		default:
			panic("no match type")
		}
		break
	case configs.MsgTypeRequest:
		switch g.Data.(type) {
		case *request.EntityInfoChangeRequest:
			pbMsgRequest = &pb.Request{
				EntityChangeRequest: g.Data.ToMessage().(*pb.EntityInfoChangeRequest),
			}
		case *request.HeartBeatRequest:
			pbMsgRequest = &pb.Request{
				HeartBeatRequest: g.Data.ToMessage().(*pb.HeartBeatRequest),
			}
		}
		break
	case configs.MsgTypeResponse:
		switch g.Data.(type) {
		case *response2.EntityInfoChangeResponse:
			pbMsgResponse = &pb.Response{
				EntityChangeResponse: g.Data.ToMessage().(*pb.EntityInfoChangeResponse),
			}
		case *response2.HeartBeatResponse:
			pbMsgResponse = &pb.Response{
				HeartBeatResponse: g.Data.ToMessage().(*pb.HeartBeatResponse),
			}
		}
		break
	default:
		panic("msg type is incorrect")
	}
	pbMsg := &pb.GMessage{
		MsgType:   pb.MSG_TYPE(g.MsgType),
		MsgCode:   pb.GAME_MSG_CODE(g.GameMsgCode),
		SessionId: g.SessionId,
		SeqId:     g.SeqId,
		Notify:    pbMsgNotify,
		Request:   pbMsgRequest,
		Response:  pbMsgResponse,
		SendTime:  tools.TIME_UTIL.NowMillis(),
	}

	return pbMsg
}

func (g *GMessage) CopyFromMessage(obj interface{}) e.Event {
	pbMsg := obj.(*pb.GMessage)
	msg := &GMessage{
		MsgType: int32(pbMsg.MsgType),
		SeqId:   pbMsg.SeqId,
	}
	msg.SetCode(int32(pbMsg.MsgCode))
	msg.SetSessionId(pbMsg.SessionId)
	msg.SetSeqId(pbMsg.SeqId)
	msg.SendTime = pbMsg.SendTime
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
	msg.Data.SetSeqId(pbMsg.SeqId)
	msg.Data.SetRoomId(g.RoomId)
	return msg
}

func (g *GMessage) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.GMessage)
	g.MsgType = int32(pbMsg.MsgType)
	g.Code = int32(pbMsg.MsgCode)
	g.SessionId = pbMsg.SessionId
	g.SeqId = pbMsg.SeqId
	g.SendTime = pbMsg.SendTime
	event := e.Manager.FetchEvent(g.GetCode())
	if pb.MSG_TYPE_NOTIFY == pbMsg.MsgType {
		g.Data = event.CopyFromMessage(pbMsg.Notify)
	}
	if pb.MSG_TYPE_REQUEST == pbMsg.MsgType {
		g.Data = event.CopyFromMessage(pbMsg.Request)
	}
	if pb.MSG_TYPE_RESPONSE == pbMsg.MsgType {
		g.Data = event.CopyFromMessage(pbMsg.Response)
	}
	g.Data.SetSessionId(pbMsg.SessionId)
	g.Data.SetRoomId(g.RoomId)
}
