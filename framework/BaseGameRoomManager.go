package framework

//import "github.com/LILILIhuahuahua/ustc_tencent_game/internal/game"

type BaseGameRoomManager struct {
	roomMap map[int64]BaseGameRoom
}

var GAME_ROOM_MANAGER = NewGameRoomManager()

func NewGameRoomManager()  *BaseGameRoomManager {
	return &BaseGameRoomManager{
		roomMap: make(map[int64]BaseGameRoom),
	}
}

func (m *BaseGameRoomManager)FetchGameRoom(id int64) BaseGameRoom {
	return m.roomMap[id]
}

func (m *BaseGameRoomManager)RegisterGameRoom(room BaseGameRoom) {
	 if nil==m.roomMap[room.GetRoomID()] {
	 	m.roomMap[room.GetRoomID()] = room
	 }
}

func (m *BaseGameRoomManager)Unicast(roomId int64, sessionId int64, buff []byte) {
	r := m.FetchGameRoom(roomId)
	s := r.FetchConnector(sessionId)
	s.SendMessage(buff)
	//m.FetchGameRoom(roomId).FetchConnector(sessionId).SendMessage(buff)
	println("hello")
}