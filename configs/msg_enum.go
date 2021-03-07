package configs

const (
	// msg消息类型
	MsgTypeNotify   int32 = 0
	MsgTypeRequest  int32 = 1
	MsgTypeResponse int32 = 2

	// game msg code
	EntityInfoChangeRequest  int32 = 0
	EntityInfoChangeResponse int32 = 1
	EntityInfoNotify         int32 = 2
	HeroQuitRequest          int32 = 3
	GameGlobalInfoNotify     int32 = 4
	TimeNotify               int32 = 5
	EnterGameNotify          int32 = 6
	EnterGameRequest         int32 = 7
	EnterGameResponse        int32 = 8
	HeroViewNotify			 int32 = 9

	// 英雄状态
	Live       int32 = 0
	Dead       int32 = 1
	Invincible int32 = 2

	// Item status
	ItemLive int32 = 0
	ItemDead int32 = 1

	// eventType的枚举
	HeroCollision int32 = 0
	ItemCollision int32 = 1
	HeroMove      int32 = 2
	HeroGrow      int32 = 3

	// entityType的枚举
	HeroType int32 = 0
	PropType int32 = 1
	FoodType int32 = 2

	// response的枚举
	Success int32 = 0
	Fail    int32 = 1

	// heroViewNotify的枚举
	Enter int32 = 0
	Leave int32 = 1
)
