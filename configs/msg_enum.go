package configs

const (
	// msg消息类型
	MsgTypeNotify   = 0
	MsgTypeRequest  = 1
	MsgTypeResponse = 2

	// game msg code
	EntityInfoChangeRequest = 0
	EntityInfoChangeResponse = 1
	EntityInfoNotify = 2
	HeroQuitRequest = 3
	GameGlobalInfoNotify = 4
	TimeNotify = 5
	EnterGameNotify = 6
	EnterGameRequest = 7
	EnterGameResponse = 8

	// 英雄状态
	Live = 0
	Dead = 1
	Invincible = 2

	// Item status
	ItemLive = 0
	ItemDead = 1

	// eventType的枚举
	HeroCollision = 0
	ItemCollision = 1
	HeroMove = 2
	HeroGrow = 3

	// entityType的枚举
	HeroType = 0
	PropType = 1
	FoodType = 2

	// response的枚举
	Success = 0
	Fail = 1
)