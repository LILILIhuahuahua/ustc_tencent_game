package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
)

type HeartBeatResponse struct {
	framework.BaseEvent
	SendTime int64
}

func NewHeartBeatResponse(time int64)  *HeartBeatResponse{
	return &HeartBeatResponse{
		SendTime: time,
	}
}

func (e *HeartBeatResponse) ToMessage() interface{} {
	return &pb.HeartBeatResponse{
		SendTime: e.SendTime,
	}
}

func (e *HeartBeatResponse) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.HeartBeatResponse)
	e.SetCode(int32(pb.GAME_MSG_CODE_HEART_BEAT_RESPONSE))
	e.SendTime = pbMsg.GetSendTime()
}

func (e *HeartBeatResponse) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Response).HeartBeatResponse
	resp := &HeartBeatResponse{
		SendTime: pbMsg.GetSendTime(),
	}
	resp.SetCode(int32(pb.GAME_MSG_CODE_HEART_BEAT_RESPONSE))
	return resp
}

func (e *HeartBeatResponse) ToGMessageBytes(seqId int32) []byte {
	resp := &pb.Response{
		HeartBeatResponse: e.ToMessage().(*pb.HeartBeatResponse),
	}
	msg := pb.GMessage{
		MsgType:  pb.MSG_TYPE_RESPONSE,
		MsgCode:  pb.GAME_MSG_CODE_HEART_BEAT_RESPONSE,
		Response: resp,
		SeqId: seqId,
		SendTime: tools.TIME_UTIL.NowMillis(),
	}
	out, _ := proto.Marshal(&msg)
	return out
}
