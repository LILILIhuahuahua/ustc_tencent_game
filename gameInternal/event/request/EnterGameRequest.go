package request

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event/info"
	"github.com/LILILIhuahuahua/ustc_tencent_game/tools"
	"github.com/golang/protobuf/proto"
)

type EnterGameRequest struct {
	framework.BaseEvent
	PlayerID   int32
	Connect    info.ConnectInfo
	PlayerName string
}

func NewEnterGameRequest(playerId int32, connectInfo info.ConnectInfo, name string) *EnterGameRequest {
	return &EnterGameRequest{
		PlayerID:   playerId,
		Connect:    connectInfo,
		PlayerName: name,
	}
}

func (e *EnterGameRequest) FromMessage(obj interface{}) {
	pbMsg := obj.(*pb.EnterGameRequest)
	e.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))
	e.PlayerID = pbMsg.GetPlayerId()
	infoMsg := pbMsg.GetClientConnectMsg()
	info := info.ConnectInfo{}
	info.FromMessage(infoMsg)
	e.Connect = info
	e.PlayerName = pbMsg.GetPlayerName()
}

func (e *EnterGameRequest) CopyFromMessage(obj interface{}) event.Event {
	pbMsg := obj.(*pb.Request).EnterGameRequest
	infoMsg := pbMsg.GetClientConnectMsg()
	info := info.ConnectInfo{}
	info.FromMessage(infoMsg)
	req := &EnterGameRequest{
		PlayerID:   pbMsg.GetPlayerId(),
		Connect:    info,
		PlayerName: pbMsg.GetPlayerName(),
	}
	req.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))
	return req
}

func (e *EnterGameRequest) ToMessage() interface{} {
	infoMsg := &pb.ConnectMsg{
		Ip:   e.Connect.Ip,
		Port: e.Connect.Port,
	}
	return &pb.EnterGameRequest{
		PlayerId:         e.PlayerID,
		ClientConnectMsg: infoMsg,
		PlayerName:       e.PlayerName,
	}
}

func (e *EnterGameRequest) ToGMessageBytes() []byte {
	req := &pb.Request{
		EnterGameRequest: e.ToMessage().(*pb.EnterGameRequest),
	}
	msg := pb.GMessage{
		MsgType:  pb.MSG_TYPE_REQUEST,
		MsgCode:  pb.GAME_MSG_CODE_ENTER_GAME_REQUEST,
		Request:  req,
		SeqId:    -1,
		SendTime: tools.TIME_UTIL.NowMillis(),
	}
	out, _ := proto.Marshal(&msg)
	return out
}
