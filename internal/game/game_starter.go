package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/configs"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/scheduler"
	"time"
)

type GameStarter struct {
	room *GameRoom
}

func NewGameStarter(addr string) *GameStarter {
	g := &GameStarter{
		room: NewGameRoom(addr),
	}
	g.init()
	return g
}

func (this *GameStarter) init() {
	//todo:加载配置

	//初始化系统组件
	GAME_ROOM_MANAGER.RegisterGameRoom(this.room)

	enterGameNotify := notify.EnterGameNotify{}
	enterGameNotify.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY))
	gameGlobalInfoNotify := notify.GameGlobalInfoNotify{}
	gameGlobalInfoNotify.SetCode(int32(pb.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY))
	entityInfochangeReq := request.EntityInfoChangeRequest{}
	entityInfochangeReq.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))
	entityInfochangeResp := response.EntityInfoChangeResponse{}
	entityInfochangeResp.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE))
	enterGameRequest := request.EnterGameRequest{}
	enterGameRequest.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))

	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY), &enterGameNotify, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST), &entityInfochangeReq, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST), &enterGameRequest, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY), &gameGlobalInfoNotify, GAME_EVENT_HANDLER)

	//todo:启动定时任务
	//定时检测客户端kcp是否可连通 每5秒检测一次
	//scheduler.NewTimer(time.Second*time.Duration(5), GAME_ROOM_MANAGER.DeleteUnavailableSession)
	//scheduler.NewTimer(time.Second*time.Duration(5), GAME_ROOM_MANAGER.DeleteDeprecatedHero)
	//定期更新小球位置，每50ms检测一次
	scheduler.NewTimer(time.Millisecond*time.Duration(50), GAME_ROOM_MANAGER.UpdateHeroPosition)
	go scheduler.Sched(configs.GlobalScheduleConfig)
}

func (this *GameStarter) Boot() {
	this.room.Serv()
}
