package ws

import "sync"

type RoomType int

const (
	TopicRoom RoomType = iota
	UserRoom
)

type Room struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Type  RoomType `json:"type"`
	Users []string `json:"-"`
}

type RoomStore interface {
	Create(roomID string, roomName string, roomType RoomType)
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

func (r *InMemoryRoomStore) Create(roomID string, roomName string, roomType RoomType) {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if !ok {
		r.Lock()
		r.rooms[roomID] = Room{
			ID:    roomID,
			Name:  roomName,
			Type:  roomType,
			Users: []string{},
		}
		r.Unlock()
	}
}

func (r *InMemoryRoomStore) Join(roomID string, userID string) bool {
	r.Lock()
	_, ok := r.rooms[roomID]
	r.Unlock()

	if ok {
		r.Lock()
		tmp := r.rooms[roomID]
		tmp.Users = append(tmp.Users, userID)
		r.rooms[roomID] = tmp
		r.Unlock()
	}
	return ok
}

func (r *InMemoryRoomStore) Leave(roomID string, userID string) {
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
		for i, cid := range r.rooms[roomID].Users {
			if userID == cid {
				tmp := r.rooms[roomID]
				tmp.Users = append(tmp.Users[:i], tmp.Users[i+1:]...)
				r.rooms[roomID] = tmp
				break // TODO: this break assumes users can join one room at a time
			}
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
		for _, cid := range room.Users {
			if cid == userID {
				rm = room
				r.Unlock()
				return rm, true
			}
		}
	}
	r.Unlock()
	return rm, false
}
