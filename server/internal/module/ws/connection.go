package ws

import (
	"sync"
)

type ConnectionStore interface {
	Store(clientID string, conn *Client)
	Load(clientID string) (conn *Client, ok bool)
	Delete(clientID string)
}

type InMemoryConnectionStore struct {
	sync.Mutex
	connections map[string]*Client
}

var _ ConnectionStore = (*InMemoryConnectionStore)(nil)

func NewInMemoryConnectionStore() *InMemoryConnectionStore {
	return &InMemoryConnectionStore{
		connections: map[string]*Client{},
	}
}

func (s *InMemoryConnectionStore) Store(clientID string, conn *Client) {
	s.Lock()
	s.connections[clientID] = conn
	s.Unlock()
}

func (s *InMemoryConnectionStore) Load(clientID string) (conn *Client, ok bool) {
	s.Lock()
	conn, ok = s.connections[clientID]
	s.Unlock()
	return conn, ok
}

func (s *InMemoryConnectionStore) Delete(clientID string) {
	s.Lock()
	delete(s.connections, clientID)
	s.Unlock()
}
