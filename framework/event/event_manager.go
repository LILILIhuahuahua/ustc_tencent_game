package event

import "sync"

type EventManager struct {
	rw         *sync.RWMutex
	handlerMap map[int32]Handler
	eventMap   map[int32]Event
}

var Manager = NewEventManager()

func NewEventManager() *EventManager {
	return &EventManager{
		rw:         &sync.RWMutex{},
		handlerMap: make(map[int32]Handler),
		eventMap:   make(map[int32]Event),
	}
}

func (e *EventManager) Register(msgCode int32, event Event, handler Handler) {
	if nil != e.handlerMap[msgCode] {
		return
	}
	if nil != handler {
		e.rw.Lock()
		e.handlerMap[msgCode] = handler
		e.rw.Unlock()
	}
	if nil != e.eventMap[msgCode] {
		return
	}
	if nil != event {
		e.rw.Lock()
		e.eventMap[msgCode] = event
		e.rw.Unlock()
	}
}

func (e *EventManager) FetchHandler(msgCode int32) Handler {
	return e.handlerMap[msgCode]
}

func (e *EventManager) FetchEvent(msgCode int32) Event {
	return e.eventMap[msgCode]
}
