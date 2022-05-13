package device

import (
	"os"
	"os/signal"

	"github.com/Masterminds/log-go"

	"kees/media-controller/helpers"
)

func (c *MediaController) ControlHandler() {
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
			log.Info("Got control command for ", state)

		case payload := <-c.Inbox:
			log.Info("Got inbound message ", payload.State)
			helpers.Dump(payload)

		case <-terminate:
			log.Info("Terminating Control Handler")
			return

		case <-interrupt:
			log.Info("Application close signal received on ControlHandler, starting teardown")
			c.Teardown()
		}
	}
}
