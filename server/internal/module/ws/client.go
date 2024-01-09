package ws

import (
	"encoding/json"
	"io"

	"github.com/gorilla/websocket"
)

type Client struct {
	*websocket.Conn
	ClientID string
}

func (c *Client) ReadJSON(v interface{}) error {
	_, r, err := c.NextReader()
	if err != nil {
		return err
	}

	err = json.NewDecoder(r).Decode(v)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}

	return err
}
