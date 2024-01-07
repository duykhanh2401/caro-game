package ws

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room, 0),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case register := <-h.Register:
			{
				if _, ok := h.Rooms[register.RoomID]; ok {
					r := h.Rooms[register.RoomID]

					if _, ok := r.Clients[register.ID]; !ok {
						r.Clients[register.ID] = register
					}
				}
			}
		case unregister := <-h.Unregister:
			{
				if _, ok := h.Rooms[unregister.RoomID]; ok {
					if _, ok := h.Rooms[unregister.RoomID].Clients[unregister.ID]; ok {
						if len(h.Rooms[unregister.RoomID].Clients) != 0 {
							h.Broadcast <- &Message{
								Content:  "user left the chat",
								RoomID:   unregister.RoomID,
								Username: unregister.Username,
							}
						}

						delete(h.Rooms[unregister.RoomID].Clients, unregister.ID)
						close(unregister.Message)
					}
				}
			}
		case message := <-h.Broadcast:
			{
				if _, ok := h.Rooms[message.RoomID]; ok {

					for _, cl := range h.Rooms[message.RoomID].Clients {
						cl.Message <- message
					}
				}
			}
		}

	}
}
