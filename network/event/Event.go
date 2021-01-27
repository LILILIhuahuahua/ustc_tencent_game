package event

type Event interface {
	GetCode() int
	SetCode(eventCode int)
	ToMessage() interface{}
	FromMessage(obj interface{})
}
