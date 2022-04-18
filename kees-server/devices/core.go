package devices

import ()

type Broker interface {
	Register(MediaController)
	Deregister(MediaController)
}

type MediaController interface {
}
