package device

import (
	"os"
	"os/signal"
	"sync"

	"github.com/Masterminds/log-go"
	"github.com/gorilla/websocket"

	"kees-client/config"
	"kees-client/helpers"
	"kees-client/messages"
)

type Client struct {
	Server  config.ServerConfig `json:"server"`
	Device  Device              `json:"device"`
	Token   string              `json:"token"`
	Auth    JWT                 `json:"auth"`
	Active  sync.WaitGroup
	Conn    *websocket.Conn `json:"conn"`
	State   string          `json:"state"`
	Outbox  chan messages.WebSocket
	Control chan messages.WebSocket
}

type Device struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Controller string `json:"controller"`
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

func NewClient(config *config.Config) *Client {
	// TODO: add custom logger on instantiation and avoid
	//       having to call c.identify on all strings

	log.Info("Creating client for " + config.Device.Name + "/" + config.Device.Version + "/" + config.Device.Controller)

	return &Client{
		Server: config.Server,
		Device: Device{
			Name:       config.Device.Name,
			Version:    config.Device.Version,
			Controller: config.Device.Controller,
		},
		Token:   config.Device.Token,
		State:   "auth",
		Outbox:  make(chan messages.WebSocket, 256),
		Control: make(chan messages.WebSocket, 128),
	}
}

func (c *Client) baseURL(scheme string) string {
	return scheme + "://" + c.Server.Host + ":" + c.Server.Port
}

func (c *Client) Run() {
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

	log.Info("Starting websocket handlers")

	c.Active.Add(1)
	go c.ReadHandler()
	go c.WriteHandler()

	c.WebSocketAuth()
	c.Active.Wait()
}

func (c *Client) Disconnect(message messages.WebSocket) {
	c.State = "disconnect"
	c.Control <- message
}

func (c *Client) WebSocketAuth() {
	data := messages.WebSocket{
		State:   "auth",
		Message: "Authenticating " + c.Device.Name,
		Data: map[string]interface{}{
			"token": c.Auth.Token,
		},
	}

	helpers.Dump(data)
	c.Outbox <- data
}

func (c *Client) ReadHandler() {
	for {
		payload := messages.WebSocket{}
		err := c.Conn.ReadJSON(&payload)
		if err != nil {
			helpers.Dump(err)
			data := messages.WebSocket{
				State:   "error",
				Message: "Invalid JSON Payload",
				Data:    map[string]interface{}{},
			}

			c.Disconnect(data)
			break
		}

		helpers.Dump(payload)

		msg := payload.State + " - " + payload.Message
		log.Info(msg)

		state := payload.State
		helpers.Dump(state)
	}
}

func (c *Client) WriteHandler() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-interrupt:
			data := messages.WebSocket{
				State:   "error",
				Message: "Application Close",
				Data:    map[string]interface{}{},
			}

			c.Disconnect(data)

		// TODO: might want to move c.Control to separate controlHandler(+goroutine)
		//       and signal to writeHandler separately
		case disconnect, ok := <-c.Control:
			helpers.Debug(ok)
			helpers.Debug(disconnect)

			// TODO: do i need SetWriteDeadline here?
			err := c.Conn.WriteJSON(disconnect)
			if err != nil {
				log.Error("WebSocket Control WriteJSON failed")
				helpers.Dump(err)
			}
			c.Active.Done()
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
