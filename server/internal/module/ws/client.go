package ws

import "github.com/gorilla/websocket"

type Client struct {
	*websocket.Conn
	ClientID string
}
