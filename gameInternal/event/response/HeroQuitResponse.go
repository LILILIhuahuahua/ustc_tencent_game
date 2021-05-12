package response

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
)

type HeroQuitResponse struct {
	framework.BaseEvent
	QuitResult bool
}

func NewHeroQuitResponse(rlt bool) *HeroQuitResponse {
	return &HeroQuitResponse{
		QuitResult: rlt,
	}
}

func (e *HeroQuitResponse) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.HeroQuitResponse)
	e.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_RESPONSE))
	e.QuitResult = pbMsg.GetQuitResult()
}

func (e *HeroQuitResponse) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Response).HeroQuitResponse
	resp := &HeroQuitResponse{
		QuitResult: pbMsg.GetQuitResult(),
	}
	resp.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_RESPONSE))
	return resp
}

func (e *HeroQuitResponse) ToMessage() interface{} {
	return &pb.HeroQuitResponse{
		QuitResult: e.QuitResult,
	}
}

func (e *HeroQuitResponse) ToGMessageBytes(seqId int32) []byte {
	resp := &pb.Response{
		HeroQuitResponse: e.ToMessage().(*pb.HeroQuitResponse),
	}
	msg := pb.GMessage{
		MsgType:  pb.MSG_TYPE_RESPONSE,
		MsgCode:  pb.GAME_MSG_CODE_HERO_QUIT_RESPONSE,
		Response: resp,
		SeqId:    seqId,
		SendTime: tools.TIME_UTIL.NowMillis(),
	}
	out, _ := proto.Marshal(&msg)
	return out
}
