package framework

type BaseGameRoomManager interface {
	FetchGameRoom(id int64) BaseGameRoom
	RegisterGameRoom(room BaseGameRoom)
	Unicast(roomId int64, sessionId int32, buff []byte)
	Braodcast(roomId int64, buff []byte)
}
