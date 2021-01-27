package event

type EventManager struct {
	handlerMap map[int32]Handler
	eventMap map[int32]Event
}

var Manager = &EventManager{
	handlerMap: make(map[int32]Handler),
	eventMap: make(map[int32]Event),
}

func (this *EventManager)NewEventManager() *EventManager{
	return &EventManager{
		handlerMap: make(map[int32]Handler),
		eventMap: make(map[int32]Event),
	}
}

func (this *EventManager)Register(msgCode int32, event Event, handler Handler) {
	if nil!=this.handlerMap[msgCode] {
		return
	}
	if nil!=handler {
		this.handlerMap[msgCode] = handler
	}
	if nil!=this.eventMap[msgCode] {
		return
	}
	if nil!=event {
		this.eventMap[msgCode] = event
	}
}

func (this *EventManager)FetchHandler(msgCode int32)Handler {
	return this.handlerMap[msgCode]
}

func (this *EventManager)FetchEvent(msgCode int32)Event {
	return this.eventMap[msgCode]
}
