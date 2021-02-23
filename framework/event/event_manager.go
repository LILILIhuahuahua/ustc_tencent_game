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

func (this *EventManager) Register(msgCode int32, event Event, handler Handler) {
	if nil != this.handlerMap[msgCode] {
		return
	}
	if nil != handler {
		this.rw.Lock()
		this.handlerMap[msgCode] = handler
		this.rw.Unlock()
	}
	if nil != this.eventMap[msgCode] {
		return
	}
	if nil != event {
		this.rw.Lock()
		this.eventMap[msgCode] = event
		this.rw.Unlock()
	}
}

func (this *EventManager) FetchHandler(msgCode int32) Handler {
	return this.handlerMap[msgCode]
}

func (this *EventManager) FetchEvent(msgCode int32) Event {
	return this.eventMap[msgCode]
}
