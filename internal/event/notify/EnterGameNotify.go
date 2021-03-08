package notify

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
)

type EnterGameNotify struct {
	framework.BaseEvent
	PlayerID int32
	Connect  info.ConnectInfo
}

func (e *EnterGameNotify) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EnterGameNotify)
	e.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY))
	e.PlayerID = pbMsg.GetPlayerId()
	infoMsg := pbMsg.GetClientConnectMsg()
	info := info.ConnectInfo{}
	info.FromMessage(infoMsg)
	e.Connect = info
}

func (e *EnterGameNotify) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Notify).EnterGameNotify
	infoMsg := pbMsg.GetClientConnectMsg()
	info := info.ConnectInfo{}
	info.FromMessage(infoMsg)
	notify := &EnterGameNotify{
		PlayerID: pbMsg.GetPlayerId(),
		Connect:  info,
	}
	notify.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY))
	return notify
}

func (e *EnterGameNotify) ToMessage() interface{} {
	infoMsg := &pb.ConnectMsg{
		Ip:   e.Connect.Ip,
		Port: e.Connect.Port,
	}
	return &pb.EnterGameNotify{
		PlayerId:         e.PlayerID,
		ClientConnectMsg: infoMsg,
	}
}
