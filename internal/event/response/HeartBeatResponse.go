package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
)

type HeartBeatResponse struct {
	framework.BaseEvent
	SendTime int64
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

func ToHeartBeatRespPBMsg(sendTime int64) *pb.HeartBeatResponse{
	return &pb.HeartBeatResponse{
		SendTime: sendTime,
	}
}
