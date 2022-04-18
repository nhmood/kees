package broker

import (
	"kees-server/devices"
	"kees-server/devices/mc"
	"kees-server/helpers"
)

type Broker struct {
	MediaControllers map[string]*mc.MediaController

	registerChan   chan *devices.MediaController
	deregisterChan chan *devices.MediaController
}

func NewBroker() *Broker {
	return &Broker{
		MediaControllers: make(map[string]*mc.MediaController),
		registerChan:     make(chan *devices.MediaController, 128),
		deregisterChan:   make(chan *devices.MediaController, 128),
	}
}

func (b *Broker) Run() {
	go b.EventHandler()

}

func (b *Broker) Register(mc *devices.MediaController) {
	b.registerChan <- mc
}

func (b *Broker) Deregister(mc *devices.MediaController) {
	b.deregisterChan <- mc
}

func (b *Broker) EventHandler() {

	for {
		select {
		case mc, ok := <-b.registerChan:
			helpers.Dump(ok)
			b.MediaControllers[mc.Info.ID] = mc
			helpers.Dump(b)

		}
	}
}
