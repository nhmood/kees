package devices

import (
	"sync"

	"github.com/Masterminds/log-go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"kees-server/helpers"
	"kees-server/messages"
)

type MediaController struct {
	Identifier string
	Info       MediaControllerInfo
	Active     sync.WaitGroup
	Conn       *websocket.Conn
	Broker     *Broker
	State      string
	Outbox     chan messages.WebSocket
	Control    chan messages.WebSocket
}

type MediaControllerInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Controller string `json:"controller"`
}

type MediaControllerState func(*MediaController, messages.WebSocket) *messages.WebSocket

var MediaControllerStates = map[string]MediaControllerState{
	"auth": (*MediaController).Auth,
}

func NewMediaController(conn *websocket.Conn, broker *Broker) *MediaController {
	id := uuid.New()
	socketID := id.String()

	// TODO: add custom logger on instantiation and avoid
	//       having to call mc.identify on all strings

	return &MediaController{
		Identifier: socketID,
		Conn:       conn,
		Broker:     broker,
		State:      "auth",
		Outbox:     make(chan messages.WebSocket, 256),
		Control:    make(chan messages.WebSocket, 128),
	}
}

func (mc *MediaController) identify(s string) string {
	return mc.Identifier + " » " + s
}

func (mc *MediaController) Run() {
	defer func() {
		log.Info(mc.identify("Session ended"))
		mc.Conn.Close()
		log.Info(mc.identify("WebSocket connection closed"))
		mc.Broker.Deregister(mc)
		log.Info(mc.identify("MediaController deregistered from broker"))
	}()

	log.Info(mc.identify("Starting MediaController WebSocket handling"))

	mc.Active.Add(1)
	go mc.ReadHandler()
	go mc.WriteHandler()
	mc.Active.Wait()
}

func (mc *MediaController) ReadHandler() {
	for {
		helpers.Debug(mc.Active)
		payload := messages.WebSocket{}
		err := mc.Conn.ReadJSON(&payload)
		if err != nil {
			helpers.Dump(err)
			data := messages.WebSocket{
				State:   "error",
				Message: "Invalid JSON Payload",
				Data:    map[string]interface{}{},
			}

			mc.Disconnect(data)
			break
		}

		msg := payload.State + " - " + payload.Message
		log.Info(mc.identify(msg))

		state := payload.State
		stateFunc := MediaControllerStates[state]
		if stateFunc == nil {
			log.Error(mc.identify("Invalid State"))
			err := messages.WebSocket{
				State:   "error",
				Message: "Invalid State",
				Data:    map[string]interface{}{},
			}
			mc.Disconnect(err)
			break
		}

		stateError := stateFunc(mc, payload)
		if stateError != nil {
			log.Error(mc.identify("Error from state:" + state + " / " + stateError.Message))

			helpers.Debug(stateError)
			mc.Disconnect(*stateError)
			return
		}
	}
}

func (mc *MediaController) WriteHandler() {
	// TODO: add ticket case for periodic heartbeat/status
	for {
		select {
		// TODO: might want to move mc.Control to separate controlHandler(+goroutine)
		//       and signal to writeHandler separately
		case disconnect, ok := <-mc.Control:
			helpers.Debug(ok)
			helpers.Debug(disconnect)

			// TODO: do i need SetWriteDeadline here?
			err := mc.Conn.WriteJSON(disconnect)
			if err != nil {
				log.Error(mc.identify("WebSocket Control WriteJSON failed"))
				helpers.Dump(err)
			}
			mc.Active.Done()
			return

		case message, ok := <-mc.Outbox:
			helpers.Debug(ok)
			helpers.Debug(message)

			// TODO: do i need SetWriteDeadline here?
			err := mc.Conn.WriteJSON(message)
			if err != nil {
				log.Error(mc.identify("WebSocket Outbox WriteJSON failed"))
				helpers.Dump(err)
			}
		}
	}
}

func (mc *MediaController) Auth(payload messages.WebSocket) *messages.WebSocket {
	log.Info(mc.identify("Auth State Received"))
	helpers.Debug(payload)

	token := payload.Data["token"]
	if token == nil {
		log.Error(mc.identify("No JWT Provided"))
		jwtError := &messages.WebSocket{
			State:   "error",
			Message: "Invalid JWT",
			Data:    map[string]interface{}{},
		}
		return jwtError
	}

	jwt, err := helpers.ValidateJWT(token.(string))
	if err != nil {
		log.Error(mc.identify("Invalid JWT from"))
		jwtError := &messages.WebSocket{
			State:   "error",
			Message: "Invalid JWT",
			Data:    map[string]interface{}{},
		}
		return jwtError
	}

	helpers.ToStruct(jwt["kees"], &mc.Info)
	log.Info(mc.identify(mc.Info.ID + " successfuly authenticated"))
	helpers.Debug(mc.Info)

	mc.Identifier = mc.Identifier + "/" + mc.Info.ID

	auth := messages.WebSocket{
		State:   "auth",
		Message: "successfully authenticated",
		Data:    helpers.ToInterface(mc.Info),
	}

	mc.Outbox <- auth

	return nil
}

func (mc *MediaController) Disconnect(message messages.WebSocket) {
	mc.State = "disconnect"
	mc.Control <- message
}
