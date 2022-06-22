package device

import (
	"github.com/Masterminds/log-go"

	"kees/controller/helpers"
	"kees/controller/messages"
)

func (c *Controller) WebSocketAuth(payload *messages.WebSocket) *messages.WebSocket {
	data := messages.WebSocket{
		State:   "auth",
		Message: "Authenticating " + c.Device.Name,
		Data: map[string]interface{}{
			"token": c.Auth.Token,
		},
	}

	helpers.Debug(data)
	c.Outbox <- data

	return nil
}

func (c *Controller) WebSocketAuthAck(payload *messages.WebSocket) *messages.WebSocket {
	log.Info("WebSocketAuth Ackd / DeviceID:" + payload.Data["id"].(string) + " name:" + payload.Data["name"].(string) + " controller:" + payload.Data["controller"].(string) + " version:" + payload.Data["version"].(string))
	helpers.Debug(payload)
	return nil
}
