package device

import (
	"bytes"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Masterminds/log-go"
	"github.com/gorilla/websocket"

	"kees-client/config"
	"kees-client/constants"
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
	c.Authenticate()



}

func (c *Client) baseURL() string {
	return "http://" + c.Server.Host + ":" + c.Server.Port
}

func (c *Client) getAuthURL() string {
	return c.baseURL() + "/ws/v1/auth"
}

func (c *Client) Authenticate() {
	log.Info("Authenticating " + c.Device.Name)
	jsonData, err := helpers.Format(c.Device)
	if err != nil {
		log.Warn("Failed to format Device info")
		os.Exit(1)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.getAuthURL(),
		bytes.NewBuffer(jsonData),
	)
	request.Header.Set("User-Agent", "kees-client/"+constants.Version)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Kees-MC-Token", c.Token)

	httpClient := http.Client{Timeout: time.Duration(5 * time.Second)}
	log.Info("Making request for: " + c.getAuthURL())
	resp, err := httpClient.Do(request)

	if err != nil {
		log.Warn("Failed to authenticate " + c.Device.Name)
		log.Error(err)
		os.Exit(1)
	}

	authResp := AuthResponse{}
	helpers.Parse(resp, &authResp)
	helpers.Debug(authResp)

	c.Device = authResp.Device
	c.Auth = authResp.JWT

	log.Info("Authentication successful - DeviceID:" + c.Device.ID)
}
