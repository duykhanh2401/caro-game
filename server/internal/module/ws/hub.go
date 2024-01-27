package ws

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	ErrBadRequest  = errors.New("Bad Request")
	ErrServerError = errors.New("Server Error")
	ErrNotFound    = errors.New("Data Not Found")
)

type Message struct {
	ID        string `json:"id"`
	UserID    string `json:"userId"`
	User      *User  `json:"user"`
	RoomID    string `json:"roomId"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Hub struct {
	Register       chan *Client
	Unregister     chan *Client
	GetRooms       chan *Client
	JoinRoom       chan *Request
	LeaveRoom      chan *Request
	SendMessage    chan *Request
	CreateRoom     chan *Request
	ChangeUserName chan *Request
	Options        *HubOptions
	connection     ConnectionStore
	room           RoomStore
	user           UserStore
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
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		GetRooms:       make(chan *Client),
		JoinRoom:       make(chan *Request),
		LeaveRoom:      make(chan *Request),
		SendMessage:    make(chan *Request),
		CreateRoom:     make(chan *Request),
		ChangeUserName: make(chan *Request),
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
		case conn := <-h.GetRooms:
			{
				h.getRooms(conn)
			}
		case req := <-h.CreateRoom:
			{
				h.createRoom(req)
			}
		case req := <-h.JoinRoom:
			{
				h.joinRoom(req)
			}
		case req := <-h.ChangeUserName:
			{
				h.changeUsername(req)
			}
		case req := <-h.SendMessage:
			{
				h.sendMessage(req)
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
				case CHANGE_USERNAME:
					{
						h.ChangeUserName <- &request
					}
				case SEND_MESSAGE:
					{
						h.SendMessage <- &request
					}
				case GET_ROOMS:
					{
						h.GetRooms <- client
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

func (h *Hub) changeUsername(req *Request) {
	fmt.Println("Change Username !!!")
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrBadRequest)
		h.unregister(conn)
		return
	}

	var username string
	{
		if tmp, ok := req.Body["username"]; ok {
			if s, ok := tmp.(string); ok {
				username = s
			} else {
				h.error(conn, ErrBadRequest)
				return
			}
		} else {
			h.error(conn, ErrBadRequest)
			return
		}
	}

	// Load user
	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrNotFound)
		return
	}

	user.Username = username
	h.user.Store(user.ID, user)

	// Inform user itself here
	res := Response{
		Body: map[string]interface{}{
			"message": "your username is changed",
			"data":    &user,
		},
		Type: ME_CHANGED_USERNAME,
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
		return
	}

	roomID := EncodeToString(6)

	var roomName string
	{
		if tmp, ok := req.Body["roomName"]; ok {
			if s, ok := tmp.(string); ok {
				roomName = s
			} else {
				h.error(conn, ErrBadRequest)
				return
			}
		} else {
			h.error(conn, ErrBadRequest)
			return
		}
	}

	roomCreated, err := h.room.Create(roomID, roomName, req.ClientID, TopicRoom)
	if err != nil {
		h.error(conn, err)
		return
	}
	fmt.Println("Create Room ID: ", roomID)

	user, ok := h.user.Load(req.ClientID)
	if !ok { // TODO: handle this more user friendly way
		user.ID = "<removed>"
		user.Username = "<removed>"
	}

	res := Response{
		Body: map[string]interface{}{
			"message": "Bạn đã tạo phòng với ID: " + roomID,
			"room":    roomCreated,
			"user":    user,
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

func (h *Hub) leaveRoom(req *Request) {
	fmt.Println("Leave Room !!!")
	conn, ok := h.connection.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrBadRequest)
		h.unregister(conn)
		return
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

	h.room.Leave(roomID, req.ClientID)

}

func (h *Hub) sendMessage(req *Request) {
	fmt.Println("Send Message !!!")
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

	room, ok := h.room.Room(roomID)
	if !ok {
		h.error(conn, ErrNotFound)
	}

	if room.Guest != req.ClientID && room.Master != req.ClientID {
		h.error(conn, errors.New("Bạn không có quyền gửi tin nhắn"))
		return
	}

	var message string
	{
		if tmp, ok := req.Body["message"]; ok {
			if s, ok := tmp.(string); ok {
				message = s
			} else {
				h.error(conn, ErrBadRequest)
				return
			}
		} else {
			h.error(conn, ErrBadRequest)
			return
		}
	}

	newMessage := Message{
		ID:        req.ID,
		UserID:    req.ClientID,
		RoomID:    roomID,
		Message:   message,
		Timestamp: time.Now().Unix() * 1000,
	}

	user, ok := h.user.Load(req.ClientID)
	if !ok {
		h.error(conn, ErrNotFound)
		return
	}

	res := Response{
		Body: map[string]interface{}{
			"data": &newMessage,
		},
		Type: ME_MESSAGE_SEND,
	}

	if err := conn.WriteJSON(res); err != nil {
		if e := h.error(conn, ErrServerError); e != nil {
			h.unregister(conn)
			// return
		}
	}

	res.Type = OTHER_MESSAGE_SEND
	newMessage.User = &user

	var c *Client
	if room.Guest != req.ClientID {
		if tmp, ok := h.connection.Load(room.Guest); ok {
			c = tmp
		}
	} else if room.Master != req.ClientID {
		if tmp, ok := h.connection.Load(room.Master); ok {
			c = tmp
		}
	}

	if err := c.WriteJSON(res); err != nil {
		if e := h.error(c, ErrServerError); e != nil {
			h.unregister(c)
			// return
		}
	}
}

func (h *Hub) getRooms(client *Client) {
	rooms := h.room.Rooms()

	res := Response{
		Body: map[string]interface{}{
			"data": rooms,
		},
		Type: ME_GET_ROOMS,
	}

	if err := client.WriteJSON(res); err != nil {
		if e := h.error(client, ErrServerError); e != nil {
			h.unregister(client)
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
			Type: OTHER_LEFT_ROOM,
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
