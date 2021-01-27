package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/handler"
)

type GameStarter struct {
	room *framework.GameRoom
}

func NewGameStarter(addr string) *GameStarter {
	g:= &GameStarter{
		room: framework.NewGameRoom(addr),
	}
	g.init()
	return g
}

func (this *GameStarter)init() {
	//todo:加载配置

	//初始化系统组件
	gameEventHandler := handler.GameEventHandler{}

	entityInfochangeReq := request.EntityInfoChangeRequest{Code: int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)}
	entityInfochangeResp := response.EntityInfoChangeResponse{Code: int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE)}

	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST),entityInfochangeReq,gameEventHandler)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE),entityInfochangeResp,gameEventHandler)

	//todo:启动定时任务
}

func (this *GameStarter)Boot()  {
	this.room.Serv()
}
