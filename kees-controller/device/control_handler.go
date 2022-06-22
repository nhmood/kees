package device

import (
	"os"
	"os/signal"

	"github.com/Masterminds/log-go"

	"kees/controller/helpers"
	"kees/controller/messages"
)

type StateHandler func(*Controller, *messages.WebSocket) *messages.WebSocket

var ControlStates = map[string]StateHandler{
	"auth": (*Controller).WebSocketAuth,
}

var InboxStates = map[string]StateHandler{
	"auth":    (*Controller).WebSocketAuthAck,
	"command": (*Controller).Command,
}

func (c *Controller) ControlHandler() {
	defer c.Active.Done()

	// create standard terminate channel to signal killing goroutine entirely
	terminate := make(chan bool, 1)
	c.Handlers["control"] = terminate

	// create channel for capturing kill signal from os
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case state := <-c.Control:
			// TODO: add state delegator to actually perform actions
			log.Info("CONTROL: ", state)
			stateFunc := ControlStates[state]
			if stateFunc == nil {
				log.Error("Unknown control command ", state)
				err := formatMessage("error", "Unknown control command", nil)
				// TODO: add optional err to Teardown()
				c.Outbox <- err
				c.Teardown()
				break
			}

			stateError := stateFunc(c, nil)
			if stateError != nil {
				log.Error("State Handling failed for ", state)
				// TODO: add optional err to Teardown()
				c.Outbox <- *stateError
				c.Teardown()
				break
			}

		case payload := <-c.Inbox:
			msg := "INBOX: " + payload.State + " - " + payload.Message
			log.Info(msg)
			helpers.Debug(payload)
			stateFunc := InboxStates[payload.State]
			if stateFunc == nil {
				log.Error("Unknown inbox command ", payload.State)
				err := formatMessage("error", "Unknown inbox command", nil)
				// TODO: add optional err to Teardown()
				c.Outbox <- err
				c.Teardown()
				break
			}

			stateError := stateFunc(c, &payload)
			if stateError != nil {
				log.Error("State Handling failed for ", payload.State)
				// TODO: add optional err to Teardown()
				c.Outbox <- *stateError
				c.Teardown()
				break
			}

		case <-terminate:
			log.Info("Terminating Control Handler")
			return

		case <-interrupt:
			log.Info("Application close signal received on ControlHandler, starting teardown")
			c.Teardown()
		}
	}
}
