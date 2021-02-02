package event
//消息接口
type Event interface {
	GetCode() int32
	SetCode(eventCode int32)
	GetSessionId() int32
	SetSessionId(sessionId int32)
	GetRoomId() int64
	SetRoomId(roomId int64)
	ToMessage() interface{}
	FromMessage(obj interface{}) //构造消息
	CopyFromMessage(obj interface{}) Event //拷贝构造消息
}
