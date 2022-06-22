package device

import (
	"github.com/Masterminds/log-go"

	"kees/controller/helpers"
	"kees/controller/messages"
)

func (c *Controller) ReadHandler() {
	defer c.Active.Done()

	// gorilla websocket does not support a channel interface so we need to
	// wait for a websocket.Conn.Close() in order to return an error from Conn.ReadJSON
	// TODO: we can probably wrap websocket.Conn.ReadJSON to support a proper channel mechanism
	for {
		payload := messages.WebSocket{}
		err := c.Conn.ReadJSON(&payload)

		// ReadHandler doesn't register a standard terminate handler
		// because Conn.ReadJSON is blocking and doesn't support a select/chan interface
		// if we get a message but we are in the teardown state, exit
		if c.State == "teardown" {
			return
		}

		if err != nil {
			log.Error(err)

			// if the socket error came from something other than a explicit
			// c.Teardown (which sets the state to "teardown"), then we had a
			// connection close/error externally, so kick off a teardown
			if c.State != "teardown" {
				c.Teardown()
			}

			return
		}

		helpers.Debug(payload)

		state := payload.State
		helpers.Debug(state)
		c.Inbox <- payload
	}
}
