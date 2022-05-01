package device

import (
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
	//       having to call mc.identify on all strings

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

func (c *Client) Run() {
	helpers.Debug(c)
	// TODO: potentially store JWT on disk and only authenticate
	//		 if expiration has passed
	//		 can use the /v1/auth/check endpoint to validate JWT
	auth := c.Authenticate()
	c.Device = auth.Device
	c.Auth = auth.JWT

	conn := c.EstablishWebSocket()
	c.Conn = conn
}

func (c *Client) baseURL(scheme string) string {
	return scheme + "://" + c.Server.Host + ":" + c.Server.Port
}

func (c *Client) getAuthURL() string {
	return c.baseURL("http") + "/ws/v1/auth"
}

func (c *Client) getWSURL() string {
	return c.baseURL("ws") + "/ws/v1/mc"
}

func (c *Client) EstablishWebSocket() *websocket.Conn {
	log.Info("Establishing websocket to: ", c.getWSURL())
	conn, _, err := websocket.DefaultDialer.Dial(c.getWSURL(), nil)
	if err != nil {
		log.Fatal("Failed to establish websocket to:", c.getWSURL(), err)
	}
	defer conn.Close()
	log.Info("Successfully established websocket")

	return conn
}
