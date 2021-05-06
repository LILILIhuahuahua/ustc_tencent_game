package configs

import "github.com/LILILIhuahuahua/ustc_tencent_game/api/proto"

const (
	// msg消息类型
	MsgTypeNotify   = int32(proto.MSG_TYPE_NOTIFY)
	MsgTypeRequest  = int32(proto.MSG_TYPE_REQUEST)
	MsgTypeResponse = int32(proto.MSG_TYPE_RESPONSE)

	// game msg code
	EntityInfoChangeRequest  = int32(proto.GAME_MSG_CODE_ENTITY_INFO_CHANGE_REQUEST)
	EntityInfoChangeResponse = int32(proto.GAME_MSG_CODE_ENTITY_INFO_CHANGE_RESPONSE)
	EntityInfoNotify         = int32(proto.GAME_MSG_CODE_ENTITY_INFO_NOTIFY)
	HeroQuitRequest          = int32(proto.GAME_MSG_CODE_HERO_QUIT_REQUEST)
	GameGlobalInfoNotify     = int32(proto.GAME_MSG_CODE_GAME_GLOBAL_INFO_NOTIFY)
	TimeNotify               = int32(proto.GAME_MSG_CODE_TIME_NOTIFY)
	EnterGameNotify          = int32(proto.GAME_MSG_CODE_ENTER_GAME_NOTIFY)
	EnterGameRequest         = int32(proto.GAME_MSG_CODE_ENTER_GAME_REQUEST)
	EnterGameResponse        = int32(proto.GAME_MSG_CODE_ENTER_GAME_RESPONSE)
	HeroViewNotify           = int32(proto.GAME_MSG_CODE_HERO_VIEW_NOTIFY)

	// 英雄状态
	HeroStatusLive = int32(proto.HERO_STATUS_LIVE)
	HeroStatusDead = int32(proto.HERO_STATUS_DEAD)

	// Item status
	ItemLive = int32(proto.ITEM_STATUS_ITEM_LIVE)
	ItemDead = int32(proto.ITEM_STATUS_ITEM_DEAD)

	// eventType的枚举
	HeroCollision = int32(proto.EVENT_TYPE_HERO_COLLISION)
	ItemCollision = int32(proto.EVENT_TYPE_ITEM_COLLISION)
	HeroMove      = int32(proto.EVENT_TYPE_HERO_MOVE)
	HeroGrow      = int32(proto.EVENT_TYPE_HERO_GROW)

	// entityType的枚举
	HeroType           = int32(proto.ENTITY_TYPE_HERO_TYPE)
	PropTypeInvincible = int32(proto.ENTITY_TYPE_PROP_TYPE_INVINCIBLE)
	PropTypeSpeedUp    = int32(proto.ENTITY_TYPE_PROP_TYPE_SPEED_UP)
	PropTypeSpeedSlow  = int32(proto.ENTITY_TYPE_PROP_TYPE_SPEED_DOWN)
	PropTypeSizeDown   = int32(proto.ENTITY_TYPE_PROP_TYPE_SIZE_DOWN)
	PropTypeFood       = int32(proto.ENTITY_TYPE_PROP_TYPE_FOOD)

	// response的枚举
	Success = int32(proto.RESULT_TYPE_SUCCESS)
	Fail    = int32(proto.RESULT_TYPE_FAIL)

	// heroViewNotify的枚举
	Enter = int32(proto.VIEW_TYPE_ENTER)
	Leave = int32(proto.VIEW_TYPE_LEAVE)
)
