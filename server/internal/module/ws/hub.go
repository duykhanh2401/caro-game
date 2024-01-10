package ws

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	ErrBadRequest  = errors.New("Bad Request")
	ErrServerError = errors.New("Server Error")
	ErrNotFound    = errors.New("Data Not Found")
)

type Hub struct {
	Register    chan *Client
	Unregister  chan *Client
	JoinRoom    chan *Request
	LeaveRoom   chan *Request
	SendMessage chan *Request
	CreateRoom  chan *Request
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
		CreateRoom:  make(chan *Request),
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
		case req := <-h.CreateRoom:
			{
				h.createRoom(req)
			}
		case req := <-h.JoinRoom:
			{
				h.joinRoom(req)
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
					if e := h.error(client, ErrBadRequest); e != nil {
						return
					}
					continue
				}

				request.ID = GetRandomID()
				request.ClientID = clientID
				fmt.Println("Request: ", request)
				switch request.Type {
				case CREATE_ROOM:
					{
						fmt.Println("Start Create Room", request)
						h.CreateRoom <- &request
					}
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
		if e := h.error(conn, ErrServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}
}

func (h *Hub) createRoom(req *Request) {
	fmt.Println("Create Room !!!")
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrBadRequest)
		h.unregister(conn)
	}

	roomID := EncodeToString(6)

	if err := h.room.Create(roomID, roomID, req.ClientID, UserRoom); err != nil {
		h.error(conn, err)
		return
	}
	fmt.Println("Create Room ID: ", roomID)
	res := Response{
		Body: map[string]interface{}{
			"message": "Bạn đã tạo phòng với ID: " + roomID,
		},
		Type: ME_CREATED_ROOM,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, ErrServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}
}

func (h *Hub) joinRoom(req *Request) {
	fmt.Println("Join Room !!!")
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrBadRequest)
		h.unregister(conn)
	}

	var roomID string
	{
		if tmp, ok := req.Body["roomID"]; ok {
			if s, ok := tmp.(string); ok {
				roomID = s
			} else {
				h.error(conn, ErrBadRequest)
				return
			}
		} else {
			h.error(conn, ErrBadRequest)
			return
		}
	}

	fmt.Println("RoomID: ", roomID)
	if ok := h.room.Join(roomID, req.ClientID); !ok {
		fmt.Println("Join Room Err")
		h.error(conn, errors.New("Phòng đã đầy"))
		return
	}
	fmt.Println("Load User")

	// Load user
	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrNotFound)
		return
	}
	fmt.Println("Load Room")

	// Load Room
	room, ok := h.room.Room(roomID)
	if !ok {
		h.error(conn, ErrNotFound)
	}
	fmt.Println("Send Response")
	res := Response{
		Body: map[string]interface{}{
			"message": "Bạn đã tham gia vào phòng : " + roomID,
			"data": map[string]interface{}{
				"room": room,
			},
		},
		Type: ME_JOINED_CHAT,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, ErrServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}

	res.Body = map[string]interface{}{
		"message": "Có người chơi mới tham gia phòng",
		"data":    user,
	}
	res.Type = OTHER_JOINED_CHAT

	c, ok := h.connection.Load(room.Master)
	if ok {
		if err := c.WriteJSON(res); err != nil {
			if e := h.error(conn, ErrServerError); e != nil {
				h.unregister(conn)
				// return
			}
		}
	}
}

func (h *Hub) unregister(conn *Client) {
	fmt.Println("Unregister: ", conn.ClientID)
	clientID := conn.ClientID
	if len(clientID) == 0 {
		h.error(conn, ErrServerError)
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
					if e := h.error(c, ErrServerError); e != nil {
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

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
