package framework

import (
	"github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	//event2 "github.com/LILILIhuahuahua/ustc_tencent_game/gameInternal/event"
)

var EVENT_HANDLER = &BaseEventHandler{}

type BaseEventHandler struct {
}

func (b BaseEventHandler) OnEvent(e event.Event) {
	if nil == e {
		return
	}
	handler := event.Manager.FetchHandler(e.GetCode())
	if nil != handler {
		handler.OnEvent(e)
	}
}

func (b BaseEventHandler) OnEventToSession(e event.Event, s event.Session) {
	if nil == e {
		return
	}
	handler := event.Manager.FetchHandler(e.GetCode())

	if nil != handler {
		handler.OnEvent(e)
	}
}
