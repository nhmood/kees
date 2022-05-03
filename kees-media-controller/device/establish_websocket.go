package device

import (
	"net/http"

	"github.com/Masterminds/log-go"
	"github.com/gorilla/websocket"

	"kees/media-controller/constants"
)

func (c *MediaController) getWSURL() string {
	return c.baseURL("ws") + "/ws/v1/mc"
}

func (c *MediaController) EstablishWebSocket() *websocket.Conn {
	log.Info("Establishing websocket to: ", c.getWSURL())

	headers := make(http.Header)
	headers.Add("User-Agent", "kees/media-controller/"+constants.Version)

	conn, _, err := websocket.DefaultDialer.Dial(c.getWSURL(), headers)
	if err != nil {
		log.Fatal("Failed to establish websocket to:", c.getWSURL(), err)
	}
	log.Info("Successfully established websocket")

	return conn
}
