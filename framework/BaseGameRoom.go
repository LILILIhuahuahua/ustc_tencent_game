package framework

type BaseGameRoom interface {
GetRoomID() int64
RegisterConnector(c *BaseSession) error
FetchConnector(sessionId int32) *BaseSession
DeleteConnector(c *BaseSession) error
BroadcastAll(buff []byte) error
Unicast(buff []byte, sessionId int64) error
Serv() error
Handle(conn *BaseSession)
}
