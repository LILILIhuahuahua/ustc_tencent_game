package framework

import (
	 "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"
	event2 "github.com/LILILIhuahuahua/ustc_tencent_game/internal/event"
)
var Handler = &BaseEventHandler{}

type BaseEventHandler struct {
}

func (b BaseEventHandler) onEvent(e event.Event) {
	if nil == e {
		return
	}
	 handler :=  event.Manager.FetchHandler(e.GetCode())
	// 二级解码
	msg := e.(event2.GMessage)
	data := msg.Data
	if nil!=handler {
		handler.OnEvent(data)
	}
}

func (b BaseEventHandler) OnEventToSession(e event.Event, s event.Session) {
	if nil == e {
		return
	}
	handler :=  event.Manager.FetchHandler(e.GetCode())

	if nil!=handler {
		handler.OnEvent(e)
	}
}