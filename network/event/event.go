package event

type event interface {
	GetCode() int
	SetCode(eventCode int)
	ToMessage() interface{}
	FromMessage(obj interface{})
}
