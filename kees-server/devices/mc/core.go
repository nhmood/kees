package mc

import (
	"github.com/gorilla/websocket"
	"sync"

	"kees-server/helpers"
	"kees-server/messages"
)

type Info struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Controller string `json:"controller"`
}

type MediaControllerState func(*MediaController, messages.WebSocket) *messages.WebSocket

var MediaControllerStates = map[string]MediaControllerState{
	"auth": (*MediaController).Auth,
}

type MediaController struct {
	Info    Info
	Active  sync.WaitGroup
	Conn    *websocket.Conn
	State   string
	Outbox  chan messages.WebSocket
	Control chan messages.WebSocket
}

func New(conn *websocket.Conn) *MediaController {
	return &MediaController{
		Conn:    conn,
		State:   "auth",
		Outbox:  make(chan messages.WebSocket, 256),
		Control: make(chan messages.WebSocket, 128),
	}
}

func (mc *MediaController) Run() {
	mc.Active.Add(1)
	go mc.ReadHandler()
	go mc.WriteHandler()
	mc.Active.Wait()
}

func (mc *MediaController) Auth(payload messages.WebSocket) *messages.WebSocket {
	helpers.Dump("Auth State")
	helpers.Dump(payload)

	token := payload.Data["token"]
	if token == nil {
		jwtError := &messages.WebSocket{
			Message: "Invalid JWT",
			Data:    map[string]interface{}{},
		}
		return jwtError

	}

	jwt, err := helpers.ValidateJWT(token.(string))
	if err != nil {
		jwtError := &messages.WebSocket{
			Message: "Invalid JWT",
			Data:    map[string]interface{}{},
		}
		return jwtError
	}

	helpers.ToStruct(jwt["kees"], &mc.Info)
	helpers.Dump(mc.Info)
	return nil
}

func (mc *MediaController) Disconnect(message messages.WebSocket) {
	mc.State = "disconnect"
	mc.Control <- message
}

func (mc *MediaController) ReadHandler() {
	for {
		helpers.Dump(mc.Active)
		payload := messages.WebSocket{}
		err := mc.Conn.ReadJSON(&payload)
		if err != nil {
			helpers.Dump(err)
			data := messages.WebSocket{
				Message: "Invalid JSON Payload",
				Data:    map[string]interface{}{},
			}

			mc.Disconnect(data)
			break
		}

		state := payload.Message
		stateFunc := MediaControllerStates[state]
		stateError := stateFunc(mc, payload)
		if stateError != nil {
			helpers.Dump(err)
			mc.Disconnect(*stateError)
			return
		}
	}
}

func (mc *MediaController) WriteHandler() {
	for {
		select {
		// TODO: might want to move mc.Control to separate controlHandler(+goroutine)
		//       and signal to writeHandler separately
		case disconnect, ok := <-mc.Control:
			helpers.Dump(ok)
			helpers.Dump(disconnect)

			// TODO: do i need SetWriteDeadline here?
			err := mc.Conn.WriteJSON(disconnect)
			mc.Active.Done()
			helpers.Dump(err)
			return

		case message, ok := <-mc.Outbox:
			helpers.Dump(ok)
			helpers.Dump(message)

			// TODO: do i need SetWriteDeadline here?
			err := mc.Conn.WriteJSON(message)
			helpers.Dump(err)
		}
	}
}
