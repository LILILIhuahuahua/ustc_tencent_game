package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
)


type GameStarter struct {
	room *GameRoom
}

func NewGameStarter(addr string) *GameStarter {
	g:= &GameStarter{
		room: NewGameRoom(addr),
	}
	g.init()
	return g
}

func (this *GameStarter)init() {
	//todo:加载配置

	//初始化系统组件
	GAME_ROOM_MANAGER.RegisterGameRoom(this.room)

	enterGameNotify:=notify.EnterGameNotify{}
	enterGameNotify.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY))
	entityInfochangeReq := request.EntityInfoChangeRequest{}
	entityInfochangeReq.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	entityInfochangeResp := response.EntityInfoChangeResponse{}
	entityInfochangeResp.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE))
	enterGameRequest:=request.EnterGameRequest{}
	enterGameRequest.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))


	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY), &enterGameNotify, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST), &entityInfochangeReq, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST), &enterGameRequest, GAME_EVENT_HANDLER)

	//todo:启动定时任务
}

func (this *GameStarter)Boot()  {
	this.room.Serv()
}
