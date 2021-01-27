package event

import (
	e "github.com/LILILIhuahuahua/ustc_tencent_game/network/event"
)
type GMessage struct {
	MsgType string
	GameMsgCode string
	SeqId int32
	Msg e.Event
}

