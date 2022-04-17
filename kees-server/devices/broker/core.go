package broker

import (
	"kees-server/devices/mc"
)

type Broker struct {
}

type broker struct {
	MediaControllers map[string]*mc.MediaController

	Register   chan *mc.MediaController
	Deregister chan *mc.MediaController
}

func (Broker) New() *broker {
	return &broker{
		MediaControllers: make(map[string]*mc.MediaController),
		Register:         make(chan *mc.MediaController),
		Deregister:       make(chan *mc.MediaController),
	}
}
