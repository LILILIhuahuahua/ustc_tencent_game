package event

type Handler interface {
	//处理消息
	OnEvent(event Event)
	OnEventToSession(event Event, s Session)
}
