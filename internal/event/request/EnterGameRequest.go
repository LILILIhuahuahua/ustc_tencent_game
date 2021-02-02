package request

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/info"
)

type EnterGameRequest struct {
	framework.BaseEvent
	PlayerID int32
	Connect info.ConnectInfo
}

func (e *EnterGameRequest)FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EnterGameRequest)
	e.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))
	e.PlayerID = pbMsg.GetPlayerId()
	infoMsg := pbMsg.GetClientConnectMsg()
	info := info.ConnectInfo{}
	info.FromMessage(infoMsg)
	e.Connect = info
}

func (e *EnterGameRequest)CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Request).EnterGameRequest
	infoMsg := pbMsg.GetClientConnectMsg()
	info := info.ConnectInfo{}
	info.FromMessage(infoMsg)
	req := &EnterGameRequest{
		PlayerID: pbMsg.GetPlayerId(),
		Connect: info,
	}
	req.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))
	return req
}

func (e *EnterGameRequest)ToMessage() interface{} {
	infoMsg:=&pb.ConnectMsg{
		Ip: e.Connect.Ip,
		Port: e.Connect.Port,
	}
	return pb.EnterGameRequest{
		PlayerId: e.PlayerID,
		ClientConnectMsg: infoMsg,
	}
}

