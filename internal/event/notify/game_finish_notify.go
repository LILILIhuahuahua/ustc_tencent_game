package notify

import (
	//pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"

	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	e "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
)

//message GameFinishNotify {
//repeated HeroRankMsg heroRankMsg = 1;
//int64 finishTime = 2;
//}

type GameFinishNotify struct {
	framework.BaseEvent
	HeroRankInfos []info.HeroRankInfo
	FinishTime    int64
}

func NewGameFinishNotify(heroRankInfos []info.HeroRankInfo, finishTime int64) *GameFinishNotify {
	return &GameFinishNotify{
		HeroRankInfos: heroRankInfos,
		FinishTime:    finishTime,
	}
}

func (notify *GameFinishNotify) FromMessage(obj interface{}) {
	//不需要做,因为这个消息只会由服务端发送给客户端，不涉及到解析
}

func (notify *GameFinishNotify) CopyFromMessage(obj interface{}) e.Event {
	//不需要做,因为这个消息只会由服务端发送给客户端，不涉及到解析
	return &GameFinishNotify{}
}

func (notify *GameFinishNotify) ToMessage() interface{} {
	pbMsg := &pb.GameFinishNotify{
		FinishTime: notify.FinishTime,
	}
	for _, heroRankInfo := range notify.HeroRankInfos {
		heroRankMsg := heroRankInfo.ToMessage().(*pb.HeroRankMsg)
		pbMsg.HeroRankMsg = append(pbMsg.HeroRankMsg, heroRankMsg)
	}
	return pbMsg
}

func (notify *GameFinishNotify) ToGMessageBytes() []byte {
	n := &pb.Notify{
		GameFinishNotify: notify.ToMessage().(*pb.GameFinishNotify),
	}
	msg := pb.GMessage{
		MsgType:  pb.MSG_TYPE_NOTIFY,
		MsgCode:  pb.GAME_MSG_CODE_GAME_FINISH_NOTIFY,
		Notify:   n,
		SendTime: tools.TIME_UTIL.NowMillis(),
	}
	out, _ := proto.Marshal(&msg)
	return out
}
