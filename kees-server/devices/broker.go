package devices

import (
	"github.com/Masterminds/log-go"
	"github.com/gorilla/websocket"

	"kees-server/helpers"
)

type Broker struct {
	MediaControllers map[string]*MediaController

	mcRegister   chan *MediaController
	mcDeregister chan *MediaController
}

func NewBroker() *Broker {
	return &Broker{
		MediaControllers: make(map[string]*MediaController),
		mcRegister:       make(chan *MediaController, 128),
		mcDeregister:     make(chan *MediaController, 128),
	}
}

func (b *Broker) Run() {
	log.Info("Starting Broker")
	go b.EventHandler()

}

func (b *Broker) RegisterMediaController(conn *websocket.Conn) *MediaController {
	log.Info("Registering new connection")
	mc := NewMediaController(conn, b)
	b.Register(mc)
	log.Info("Registered MediaController:" + mc.Identifier)

	return mc
}

func (b *Broker) Register(mc *MediaController) {
	b.mcRegister <- mc
}

func (b *Broker) Deregister(mc *MediaController) {
	b.mcDeregister <- mc
}

func (b *Broker) EventHandler() {
	for {
		select {
		case mc, ok := <-b.mcRegister:
			log.Info("Registration event for " + mc.Identifier)
			helpers.Debug(ok)
			b.MediaControllers[mc.Info.ID] = mc
			helpers.Debug(b)

		case mc, ok := <-b.mcDeregister:
			log.Info("Deregistration event for " + mc.Identifier)
			helpers.Debug(ok)
			delete(b.MediaControllers, mc.Info.ID)
			helpers.Debug(b)
		}
	}
}
