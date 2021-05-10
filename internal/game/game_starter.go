package game

import (
	pb "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/notify"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/request"
	"github.com/LILILIhuahuahua/ustc_tencent_game/internal/event/response"
	"log"
)

type GameStarter struct {
	//room *GameRoom
	roomManager *GameRoomManager
}

var GAME_ROOM_MANAGER *GameRoomManager = NewGameRoomManager()

func NewGameStarter(addr string) *GameStarter {
	g := &GameStarter{
		//room: NewGameRoom(addr),
		//roomManager: NewGameRoomManager(),
	}
	g.Init()
	return g
}

func (g *GameStarter) Init() {
	//todo:加载配置

	//初始化系统组件
	//GAME_ROOM_MANAGER.RegisterGameRoom(this.room)
	log.Println("[GameStarter]初始化系统组件！")
	enterGameNotify := notify.EnterGameNotify{}
	enterGameNotify.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY))

	gameGlobalInfoNotify := notify.GameGlobalInfoNotify{}
	gameGlobalInfoNotify.SetCode(int32(pb.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY))

	entityInfochangeReq := request.EntityInfoChangeRequest{}
	entityInfochangeReq.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST))

	entityInfochangeNotify := notify.EntityInfoChangeNotify{}
	entityInfochangeNotify.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_NOTIFY))

	entityInfochangeNotify1 := notify.EntityInfoChangeNotify{}
	entityInfochangeNotify1.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_NOTIFY))

	entityInfochangeResp := response.EntityInfoChangeResponse{}
	entityInfochangeResp.SetCode(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE))

	enterGameRequest := request.EnterGameRequest{}
	enterGameRequest.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST))

	enterGameResponse := response.EnterGameResponse{}
	enterGameResponse.SetCode(int32(pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE))

	heartBeatRequest := request.HeartBeatRequest{}
	heartBeatRequest.SetCode(int32(pb.GAME_MSG_CODE_HEART_BEAT_REQUEST))

	heroQuitRequest := request.HeroQuitRequest{}
	heroQuitRequest.SetCode(int32(pb.GAME_MSG_CODE_HERO_QUIT_REQUEST))

	gameFinishNotify := notify.GameFinishNotify{}
	gameFinishNotify.SetCode(int32(pb.GAME_MSG_CODE_GAME_FINISH_NOTIFY))

	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_NOTIFY), &enterGameNotify, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST), &entityInfochangeReq, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_CHANGE_NOTIFY), &entityInfochangeNotify, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTITY_INFO_NOTIFY), &entityInfochangeNotify1, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_REQUEST), &enterGameRequest, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_ENTER_GAME_RESPONSE), &enterGameResponse, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY), &gameGlobalInfoNotify, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_HEART_BEAT_REQUEST), &heartBeatRequest, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_HERO_QUIT_REQUEST), &heroQuitRequest, GAME_EVENT_HANDLER)
	event.Manager.Register(int32(pb.GAME_MSG_CODE_GAME_FINISH_NOTIFY), &gameFinishNotify, GAME_EVENT_HANDLER)

	//todo:启动定时任务
	//定时检测客户端kcp是否可连通 每5秒检测一次
	//scheduler.NewTimer(time.Second*time.Duration(5), GAME_ROOM_MANAGER.DeleteUnavailableSession)
	//scheduler.NewTimer(time.Second*time.Duration(5), GAME_ROOM_MANAGER.DeleteDeprecatedHero)
	//定期更新小球位置，每50ms检测一次
	//scheduler.NewTimer(time.Millisecond*time.Duration(50), GAME_ROOM_MANAGER.UpdateHeroPositionAndStatus)
	//go scheduler.Sched(configs.GlobalScheduleConfig)
}

func (this *GameStarter) Boot() {
	//this.room.Serv()

	GAME_ROOM_MANAGER.Serv()
}
