package device

import (
	"github.com/Masterminds/log-go"

	"kees/controller/helpers"
)

func (c *Controller) WriteHandler() {
	defer c.Active.Done()

	// create standard terminate channel to signal killing goroutine entirely
	terminate := make(chan bool, 1)
	c.Handlers["write"] = terminate

	for {
		select {
		case <-terminate:
			disconnect := formatMessage("close", "websocket connection terminated", nil)
			// TODO: do i need SetWriteDeadline here?
			err := c.Conn.WriteJSON(disconnect)
			if err != nil {
				log.Error("WebSocket Control WriteJSON failed")
				helpers.Dump(err)
			}
			log.Info("Terminating Write Handler")
			return

		case message, ok := <-c.Outbox:
			helpers.Debug(ok)
			helpers.Debug(message)

			// TODO: do i need SetWriteDeadline here?
			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Error("WebSocket Outbox WriteJSON failed")
				helpers.Dump(err)
			}
		}
	}
}
