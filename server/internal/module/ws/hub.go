package ws

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Register    chan *Client
	Unregister  chan *Client
	JoinRoom    chan *Request
	LeaveRoom   chan *Request
	SendMessage chan *Request
	Options     *HubOptions
	connection  ConnectionStore
	room        RoomStore
	user        UserStore
}

type HubOptions struct {
	MaxSavedMessage    int
	MaxReturnedMessage int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) Defaults() {
	h.Options = &HubOptions{
		MaxSavedMessage:    500,
		MaxReturnedMessage: 20,
	}
	h.connection = NewInMemoryConnectionStore()
	h.user = NewInMemoryUserStore()
	h.room = NewInMemoryRoomStore()
}

func NewHub() *Hub {
	return &Hub{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		JoinRoom:    make(chan *Request),
		LeaveRoom:   make(chan *Request),
		SendMessage: make(chan *Request),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			{
				h.register(conn)
			}
		case conn := <-h.Unregister:
			{
				h.unregister(conn)
			}
		}

	}
}

func (h *Hub) Handler() gin.HandlerFunc {
	fmt.Println("Start Handler")
	return func(ctx *gin.Context) {

		if websocket.IsWebSocketUpgrade(ctx.Request) {
			conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			clientID := GetRandomID()
			client := &Client{Conn: conn, ClientID: clientID}

			h.Register <- client

			defer func() {
				h.Unregister <- client

				if err := conn.Close(); err != nil {
					fmt.Printf("%#v\n", err)
				}
			}()

			for {
				var request Request
				if err := client.ReadJSON(&request); err != nil {
					fmt.Println("1", err)
					if e := h.error(client, errors.New("Bad Request")); e != nil {
						return
					}
					continue
				}

				request.ID = GetRandomID()
				request.ClientID = clientID

				switch request.Type {
				case JOIN_ROOM:
					{
						h.JoinRoom <- &request
					}
				}
			}
		}
	}
}

func (h *Hub) register(conn *Client) {
	fmt.Println("Register !!!!!!!!!", conn.ClientID)
	clientID := conn.ClientID
	if len(clientID) == 0 {
		h.unregister(conn)
		return
	}

	user := User{
		ID:       clientID,
		Username: clientID,
	}

	h.connection.Store(clientID, conn)
	h.user.Store(clientID, user)

	res := Response{
		Body: map[string]interface{}{
			"message": "connection successful",
			"data":    &user,
		},
		Type: CONNECTED,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, errors.New("Server Error")); e != nil {
			h.unregister(conn)
			// return
		}
	}
}

func (h *Hub) createRoom(conn *Client) {

}

func (h *Hub) unregister(conn *Client) {
	fmt.Println("Unregister: ", conn.ClientID)
	clientID := conn.ClientID
	if len(clientID) == 0 {
		h.error(conn, errors.New("Server Error"))
		return
	}

	user, ok := h.user.Load(clientID)
	if !ok {
		user.ID = "<removed>"
		user.Username = "<removed>"
	}

	// If user joined a chat room than get room id
	var roomID string
	{
		if room, ok := h.room.UserJoinedTo(user.ID); ok {
			roomID = room.ID
		}
	}

	h.room.Leave(roomID, clientID) // TODO: find a better way to handle unknown roomId situation
	// Delete connection
	h.connection.Delete(clientID)
	// Delete user
	h.user.Delete(clientID)

	// If user is removed than cannot inform who left the chat
	if user.ID == "<removed>" {
		return
	}

	// If user joined a chat room than get user ids in that chat
	var userIDs []string
	{
		if room, ok := h.room.Room(roomID); ok {
			userIDs = room.Users
		}
	}

	if len(userIDs) > 0 {
		res := Response{
			Body: map[string]interface{}{
				"message": "a user lost connection",
				"data":    &user,
			},
			Type: OTHER_LEFT_CHAT,
		}

		for _, userID := range userIDs {
			if c, ok := h.connection.Load(userID); ok {
				if userID == user.ID {
					continue // pass user itself
				}

				if err := c.WriteJSON(res); err != nil {
					if e := h.error(c, errors.New("Server Error")); e != nil {
						h.unregister(c)
						// return
						continue
					}
				}
			}
		}
	}
}

func (h *Hub) error(conn *Client, err error) error {
	res := Response{
		Error: map[string]interface{}{
			"message": err.Error(),
		},
		Type: ERROR,
	}
	return conn.WriteJSON(res)
}
