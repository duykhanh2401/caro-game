package ws

import (
	"errors"
	"fmt"
	"sync"
)

type RoomType int

const (
	TopicRoom RoomType = iota
	UserRoom
)

type Room struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Type         RoomType    `json:"type"`
	Users        []string    `json:"-"`
	Master       string      `json:"roomMaster"`
	MasterWin    int         `json:"roomMasterWin"`
	Guest        string      `json:"guest"`
	GuestWin     int         `json:"guestWin"`
	MasterFirst  bool        `json:"roomMasterFirst"`
	IsMasterTurn bool        `json:"isMasterTurn"`
	DataCaro     [20][20]int `json:"-"`
}

type RoomStore interface {
	Create(roomID string, roomName string, userCreate string, roomType RoomType) error
	Join(roomID string, userID string) bool
	Leave(roomID string, userID string)
	Users(roomID string) []string
	Room(roomID string) (room Room, ok bool)
	Rooms(includeUserRoom ...bool) []Room
	UserJoinedTo(userID string) (room Room, ok bool)
}

func NewInMemoryRoomStore() *InMemoryRoomStore {
	return &InMemoryRoomStore{
		rooms: map[string]Room{},
	}
}

type InMemoryRoomStore struct {
	sync.Mutex
	rooms map[string]Room
}

var _ RoomStore = (*InMemoryRoomStore)(nil)

func (r *InMemoryRoomStore) Create(roomID string, roomName string, userCreate string, roomType RoomType) error {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if !ok {
		r.Lock()
		r.rooms[roomID] = Room{
			ID:           roomID,
			Name:         roomName,
			Type:         roomType,
			Master:       userCreate,
			MasterFirst:  true,
			IsMasterTurn: true,
		}
		r.Unlock()
	} else {
		return errors.New("Phòng đã tồn tại")
	}

	return nil
}

func (r *InMemoryRoomStore) ResetDataCaro(roomID string) {

	r.Lock()
	room, ok := r.rooms[roomID]
	r.Unlock()
	if ok {
		for y := 0; y < 20; y++ {
			for x := 0; x < 20; x++ {
				room.DataCaro[x][y] = 0
			}
		}
	}

}

func (r *InMemoryRoomStore) Join(roomID string, userID string) bool {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()
	fmt.Println(ok)
	if ok {
		r.Lock()
		tmp := r.rooms[roomID]
		if tmp.Guest == "" {
			tmp.Guest = userID
			tmp.GuestWin = 0
			tmp.MasterWin = 0
			tmp.IsMasterTurn = true
			r.rooms[roomID] = tmp
		} else {
			return false
		}
		r.Unlock()
		r.ResetDataCaro(roomID)
	}
	return ok
}

func (r *InMemoryRoomStore) Leave(roomID string, userID string) {
	fmt.Println("User ID Leave Room: ", roomID, userID)
	if roomID == "" { // TODO: find a better way to handle unknown roomId situation
		r.Lock()
		for id, room := range r.rooms {
			for i, cid := range room.Users {
				if userID == cid {
					tmp := r.rooms[id]
					tmp.Users = append(tmp.Users[:i], tmp.Users[i+1:]...)
					r.rooms[id] = tmp
					break // TODO: this break assumes users can join one room at a time
				}
			}
		}
		r.Unlock()
	} else {
		r.Lock()

		room, ok := r.rooms[roomID]
		if ok {
			room.GuestWin = 0
			room.MasterWin = 0
			room.MasterFirst = true
			room.IsMasterTurn = true

			if room.Guest == userID {
				room.Guest = ""
			} else if room.Master == userID && room.Guest == "" {
				delete(r.rooms, roomID)
			} else if room.Master == userID && room.Guest != "" {
				room.Master = room.Guest
				room.Guest = ""
			}

			r.rooms[roomID] = room
		}

		r.Unlock()
	}
}

func (r *InMemoryRoomStore) Users(roomID string) []string {
	r.Lock()
	room := r.rooms[roomID]
	r.Unlock()

	return room.Users
}

func (r *InMemoryRoomStore) Room(roomID string) (room Room, ok bool) {
	r.Lock()
	room, ok = r.rooms[roomID]
	r.Unlock()
	return room, ok
}

func (r *InMemoryRoomStore) Rooms(includeUserRoom ...bool) []Room {
	incAll := false
	if len := len(includeUserRoom); len > 0 && includeUserRoom[0] {
		incAll = true
	}
	var rooms []Room
	r.Lock()
	for _, room := range r.rooms {
		if room.Type == UserRoom && !incAll {
			continue
		}

		rooms = append(rooms, room)
	}
	r.Unlock()
	return rooms
}

func (r *InMemoryRoomStore) UserJoinedTo(userID string) (room Room, ok bool) {
	var rm Room
	r.Lock()
	for _, room := range r.rooms {
		if room.Guest == userID || room.Master == userID {
			rm = room
			r.Unlock()
			return rm, true
		}
	}
	r.Unlock()
	return rm, false
}
