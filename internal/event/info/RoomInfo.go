package info

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework"
)

type RoomInfo struct {
	framework.BaseEvent
	ID           int32
	Status	     int32
	CreateTime   int64
	PlayerCount  int32
	HighestScore int32
}
