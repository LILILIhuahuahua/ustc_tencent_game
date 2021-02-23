package game
import "fmt"

type GameRoomManager struct {
	roomMap map[int64]*GameRoom
}

var GAME_ROOM_MANAGER = NewGameRoomManager()

func NewGameRoomManager() *GameRoomManager {
	return &GameRoomManager{
		roomMap: make(map[int64]*GameRoom),
	}
}

func (m *GameRoomManager) FetchGameRoom(id int64) *GameRoom {
	return m.roomMap[id]
}

func (m *GameRoomManager) RegisterGameRoom(room *GameRoom) {
	if nil == m.roomMap[room.GetRoomID()] {
		m.roomMap[room.GetRoomID()] = room
	}
}

func (m *GameRoomManager) Unicast(roomId int64, sessionId int32, buff []byte) {
	r := m.FetchGameRoom(roomId)
	s := r.FetchConnector(sessionId)
	s.SendMessage(buff)
	//m.FetchGameRoom(roomId).FetchConnector(sessionId).SendMessage(buff)
	println("hello")
}

func (m *GameRoomManager) Braodcast(roomId int64, buff []byte) {
	r := m.FetchGameRoom(roomId)
	r.BroadcastAll(buff)
	//m.FetchGameRoom(roomId).FetchConnector(sessionId).SendMessage(buff)
}

func (m * GameRoomManager) DeleteUnavailableSession() {
	for _, room := range m.roomMap {
		err := room.DeleteUnavailableSession()
		if err != nil {
			fmt.Println("清理不可用session的时候发生了error: ", err.Error())
		}
	}
}
