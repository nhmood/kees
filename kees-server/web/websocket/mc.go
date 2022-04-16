package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"

	"kees-server/helpers"
)

type WebSocketMessage struct {
	Message string
	Data    map[string]interface{}
}

type MediaController struct {
	Info    MediaControllerInfo
	Active  sync.WaitGroup
	Conn    *websocket.Conn
	State   string
	Outbox  chan WebSocketMessage
	Control chan WebSocketMessage
}

type MediaControllerState func(*MediaController, WebSocketMessage) *WebSocketMessage

var MediaControllerStates = map[string]MediaControllerState{
	"auth": (*MediaController).Auth,
}

func MediaControllerV1(w http.ResponseWriter, r *http.Request) {
	// establish websocket connection
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.Dump("Failed to upgrade ws")
		helpers.Dump(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer c.Close()

	mediaController := MediaController{
		Conn:    c,
		State:   "auth",
		Outbox:  make(chan WebSocketMessage, 256),
		Control: make(chan WebSocketMessage, 128),
	}
	helpers.Dump(mediaController)

	mediaController.Active.Add(1)
	go mediaController.readHandler()
	go mediaController.writeHandler()

	mediaController.Active.Wait()
	helpers.Dump("Closing")
}

func (mc *MediaController) Disconnect(message WebSocketMessage) {
	mc.State = "disconnect"
	mc.Control <- message
}

func (mc *MediaController) readHandler() {
	for {
		helpers.Dump(mc.Active)
		payload := WebSocketMessage{}
		err := mc.Conn.ReadJSON(&payload)
		if err != nil {
			helpers.Dump(err)
			data := WebSocketMessage{
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

func (mc *MediaController) writeHandler() {
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

func (mc *MediaController) Auth(payload WebSocketMessage) *WebSocketMessage {
	helpers.Dump("Auth State")
	helpers.Dump(payload)

	token := payload.Data["token"]
	if token == nil {
		jwtError := &WebSocketMessage{
			Message: "Invalid JWT",
			Data:    map[string]interface{}{},
		}
		return jwtError

	}

	jwt, err := helpers.ValidateJWT(token.(string))
	if err != nil {
		jwtError := &WebSocketMessage{
			Message: "Invalid JWT",
			Data:    map[string]interface{}{},
		}
		return jwtError
	}

	helpers.ToStruct(jwt["kees"], &mc.Info)
	helpers.Dump(mc.Info)
	return nil
}
