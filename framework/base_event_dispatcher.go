package framework

import "github.com/LILILIhuahuahua/ustc_tencent_game/framework/event"

type BaseEventDispatcher struct {
	maxEventSize int32
	EventQueue *event.EventRingQueue
}

func NewBaseEventDispatcher(maxEventSize int32) BaseEventDispatcher{
	return BaseEventDispatcher{
		maxEventSize: maxEventSize,
		EventQueue: event.NewEventRingQueue(maxEventSize),
	}
}

func (b BaseEventDispatcher) FireEvent(e event.Event) {
	//EVENT_HANDLER.OnEvent(e)
	b.EventQueue.Push(e)
}

func (b BaseEventDispatcher) GetEventQueue() *event.EventRingQueue  {
	return b.EventQueue
}

func (b BaseEventDispatcher) FireEventToSession(e event.Event, s event.Session) {

}

func (b BaseEventDispatcher) Close() {

}
