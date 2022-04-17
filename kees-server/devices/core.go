package devices

import (
	"kees-server/devices/broker"
	"kees-server/devices/mc"
)

type Broker broker.Broker

type MC struct {
	Source mc.MediaController
}
