package device

import (
	"github.com/Masterminds/log-go"
	"github.com/gorilla/websocket"
)

func (c *Client) getWSURL() string {
	return c.baseURL("ws") + "/ws/v1/mc"
}

func (c *Client) EstablishWebSocket() *websocket.Conn {
	log.Info("Establishing websocket to: ", c.getWSURL())
	conn, _, err := websocket.DefaultDialer.Dial(c.getWSURL(), nil)
	if err != nil {
		log.Fatal("Failed to establish websocket to:", c.getWSURL(), err)
	}
	log.Info("Successfully established websocket")

	return conn
}
