// myapp/internal/chat/roomManager.go
package chat

import (
	"fmt"
	"sync"
)

var (
	roomOnce    sync.Once
	roomManager *RoomManager
)

// Question: Do we Actually need map for rooms?
// Can we just use slice for rooms?
// -> answer: using slice would make algorithm harder to implement
// considering the fact that removing a room from the slice would require
// But still using integer instead of string for roomID might be better?
type RoomManager struct {
	rooms map[string]*Room
}

func GetRoomManager() *RoomManager {
	roomOnce.Do(func() {
		roomManager = &RoomManager{
			rooms: make(map[string]*Room),
		}
	})

	return roomManager
}

func (rm *RoomManager) CheckRoom(roomID string) bool {
	_, ok := rm.rooms[roomID]
	return ok
}

func (rm *RoomManager) GetRoom(roomID string) *Room {
	if !rm.CheckRoom(roomID) {
		fmt.Println("Room with roomID", roomID, "does not exist")
		return nil
	}

	return rm.rooms[roomID]
}

func (rm *RoomManager) AddRoom(room *Room) {
	if rm.CheckRoom(room.roomID) {
		fmt.Println("Room with roomID", room.roomID, "already exists")
		return
	}
	rm.rooms[room.roomID] = room
}

func (rm *RoomManager) RemoveRoom(roomID string) {
	if !rm.CheckRoom(roomID) {
		fmt.Println("Room with roomID", roomID, "does not exist")
		return
	}
	delete(rm.rooms, roomID)
}

// Question: wouldn't it be better to just return room pointers?
func (rm *RoomManager) GetRoomIDs() []string {
	roomIDs := make([]string, 0, len(rm.rooms))
	for roomID := range rm.rooms {
		roomIDs = append(roomIDs, roomID)
	}
	return roomIDs
}
