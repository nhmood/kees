package device

import (
	"sync"

	"github.com/Masterminds/log-go"
	"github.com/gorilla/websocket"

	"kees/controller/config"
	"kees/controller/helpers"
	"kees/controller/messages"
)

type Controller struct {
	Server   config.ServerConfig `json:"server"`
	Device   Device              `json:"device"`
	Token    string              `json:"token"`
	Auth     JWT                 `json:"auth"`
	Active   sync.WaitGroup
	Conn     *websocket.Conn `json:"conn"`
	State    string          `json:"state"`
	Handlers map[string]HandlerClose
	Outbox   chan messages.WebSocket
	Inbox    chan messages.WebSocket
	Control  chan string
}

type HandlerClose chan bool

type Device struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Controller   string            `json:"controller"`
	Capabilities map[string]string `json:"capabilities"`
}

type JWT struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type AuthResponse struct {
	Message string `json:"message"`
	Device  Device `json:"device"`
	JWT     JWT    `json:"jwt"`
}

func NewController(config *config.Config) *Controller {
	// TODO: add custom logger on instantiation and avoid
	//       having to call c.identify on all strings

	log.Info("Creating controller for " + config.Device.Name + "/" + config.Device.Version + "/" + config.Device.Controller)

	return &Controller{
		Server: config.Server,
		Device: Device{
			Name:         config.Device.Name,
			Version:      config.Device.Version,
			Controller:   config.Device.Controller,
			Capabilities: config.Capabilities,
		},
		Token:    config.Device.Token,
		State:    "auth",
		Outbox:   make(chan messages.WebSocket, 256),
		Inbox:    make(chan messages.WebSocket, 128),
		Control:  make(chan string, 128),
		Handlers: make(map[string]HandlerClose),
	}
}

func (c *Controller) baseURL(scheme string) string {
	return scheme + "://" + c.Server.Host + ":" + c.Server.Port
}

func formatMessage(state string, message string, data *map[string]interface{}) messages.WebSocket {
	if data == nil {
		data = &map[string]interface{}{}
	}

	return messages.WebSocket{
		State:   state,
		Message: message,
		Data:    *data,
	}
}

func (c *Controller) Run() {
	defer func() {
		log.Info("Session ended")
		c.Conn.Close()
	}()

	helpers.Debug(c)
	// TODO: potentially store JWT on disk and only authenticate
	//		 if expiration has passed
	//		 can use the /v1/auth/check endpoint to validate JWT
	auth := c.Authenticate()
	c.Device = auth.Device
	c.Auth = auth.JWT

	conn := c.EstablishWebSocket()
	c.Conn = conn

	c.StartHandlers()
	c.Control <- "auth"
	c.Active.Wait()
}

// NOTE: waitgroup.Add needs to be issued outside of the goroutine otherwise
//	 there is a race on the goroutine start and the waitgroup.Wait
//	 this results in the application closing when it should really be waiting on the handlers
func (c *Controller) StartHandlers() {
	go c.ReadHandler()
	c.Active.Add(1)

	go c.WriteHandler()
	c.Active.Add(1)

	go c.ControlHandler()
	c.Active.Add(1)
}

func (c *Controller) Teardown() {
	c.State = "teardown"
	log.Info("Tearing down Controller")

	// ReadHandler doesn't register a handler because Conn.ReadJSON
	// is blocking and doesn't support a select/chan interface
	// we need to just close the Conn for it to terminate
	// for everything else with a handler, push a terminate message
	// down on the registered handler chan
	for handler, terminateChan := range c.Handlers {
		log.Info("Pushing terminate to " + handler)
		terminateChan <- true
	}
}
