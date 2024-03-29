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
	Master       string      `json:"master"`
	MasterWin    int         `json:"masterWin"`
	MasterFirst  bool        `json:"masterFirst"`
	Guest        string      `json:"guest"`
	GuestWin     int         `json:"guestWin"`
	IsMasterTurn bool        `json:"isMasterTurn"`
	GuestReady   bool        `json:"guestReady"`
	MasterReady  bool        `json:"masterReady"`
	DataCaro     [221]string `json:"data"`
}

type RoomStore interface {
	Create(roomID string, roomName string, userCreate string, roomType RoomType) (Room, error)
	Join(roomID string, userID string) bool
	Leave(roomID string, userID string)
	Users(roomID string) []string
	Room(roomID string) (room Room, ok bool)
	Rooms(includeUserRoom ...bool) []Room
	UserJoinedTo(userID string) (room Room, ok bool)
	GuestReady(roomID string, isReady bool) bool
	MasterReady(roomID string, isReady bool) bool
	HandleGame(roomID string, isXTurn bool, index int32) (*Room, bool)
	HandleWinGame(roomID string, isMasterWin bool) (*Room, bool)
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

func (r *InMemoryRoomStore) Create(roomID string, roomName string, userCreate string, roomType RoomType) (Room, error) {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	room := Room{
		ID:           roomID,
		Name:         roomName,
		Type:         roomType,
		Master:       userCreate,
		MasterFirst:  true,
		IsMasterTurn: true,
	}

	if !ok {
		r.Lock()
		r.rooms[roomID] = room
		r.Unlock()
	} else {
		return Room{}, errors.New("Phòng đã tồn tại")
	}

	return room, nil
}

func (r *InMemoryRoomStore) HandleGame(roomID string, isXTurn bool, index int32) (*Room, bool) {
	r.Lock()
	defer r.Unlock()
	room, ok := r.rooms[roomID]

	if ok {
		fmt.Println(room.DataCaro[index])
		fmt.Println(room.DataCaro[index] == "")
		if room.DataCaro[index] == "" {
			if isXTurn {
				room.DataCaro[index] = "x"
			} else {
				room.DataCaro[index] = "o"
			}

			room.IsMasterTurn = !room.IsMasterTurn

			r.rooms[roomID] = room
			return &room, true
		}

	}

	return nil, false
}

func (r *InMemoryRoomStore) HandleWinGame(roomID string, isMasterWin bool) (*Room, bool) {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if ok {
		r.ResetDataCaro(roomID)
		r.Lock()

		room := r.rooms[roomID]

		if isMasterWin {
			room.MasterWin = room.MasterWin + 1
		} else {
			room.GuestWin = room.GuestWin + 1
		}

		room.IsMasterTurn = !room.MasterFirst
		room.MasterFirst = !room.MasterFirst

		room.MasterReady = false
		room.GuestReady = false

		r.rooms[roomID] = room
		r.Unlock()

		return &room, ok
	}

	return nil, ok
}

func (r *InMemoryRoomStore) ResetDataCaro(roomID string) {

	r.Lock()
	room, ok := r.rooms[roomID]
	if ok {
		for x := 0; x < 221; x++ {
			room.DataCaro[x] = ""
		}

		r.rooms[roomID] = room
	}
	r.Unlock()

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

func (r *InMemoryRoomStore) GuestReady(roomID string, isReady bool) bool {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if ok {
		r.Lock()
		room := r.rooms[roomID]
		room.GuestReady = isReady
		r.rooms[roomID] = room
		r.Unlock()
	}

	return ok
}

func (r *InMemoryRoomStore) MasterReady(roomID string, isReady bool) bool {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if ok {
		r.Lock()
		room := r.rooms[roomID]
		room.MasterReady = isReady
		r.rooms[roomID] = room
		r.Unlock()
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
			room.GuestReady = false
			room.MasterReady = false
			if room.Guest == userID {
				room.Guest = ""
			} else if room.Master == userID && room.Guest == "" {
				delete(r.rooms, roomID)
				r.Unlock()
				return
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
		if (room.Type == UserRoom && !incAll) || (room.Guest != "" && room.Master != "") {
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
