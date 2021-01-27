package event

type EventDispatcher interface {
	 //分发消息
	 FireEvent(event Event)
	 FireEventToSession(event Event, s Session)
	 //关闭消息分发器
	 Close()
}
