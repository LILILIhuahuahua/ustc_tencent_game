package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/handler"
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
	framework.GAME_ROOM_MANAGER.RegisterGameRoom(this.room)

	enterGameNotify:=notify.EnterGameNotify{}
	enterGameNotify.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY))
	entityInfochangeReq := request.EntityInfoChangeRequest{}
	entityInfochangeReq.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	entityInfochangeResp := response.EntityInfoChangeResponse{}
	entityInfochangeResp.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE))


	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY), &enterGameNotify,handler.GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST), &entityInfochangeReq,handler.GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE), &entityInfochangeResp,handler.GAME_EVENT_HANDLER)

	//todo:启动定时任务
}

func (this *GameStarter)Boot()  {
	this.room.Serv()
}
