package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"kees/server/devices"
	"kees/server/helpers"
	"kees/server/models"
	"kees/server/web/middlewares"
	"kees/server/web/responses"
)

var upgrader = websocket.Upgrader{}
var broker *devices.Broker

type DeviceUpdate struct {
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

type AuthResponse struct {
	Message string        `json:"message"`
	Device  models.Device `json:"device"`
	JWT     JWTResponse   `json:"jwt"`
}

type JWTResponse struct {
	ExpiresIn int64  `json:"expires_in"`
	Token     string `json:"token"`
}

func Configure(router *mux.Router, path string, b *devices.Broker) {
	broker = b

	ws := router.PathPrefix(path).Subrouter()
	ws.Use(middlewares.AddJSONHeader)

	ws.HandleFunc("/", WebsocketInfo).Methods("GET")
	ws.HandleFunc("/v1/auth", WebsocketAuthV1).Methods("POST")
	ws.HandleFunc("/v1/mc", MediaControllerV1).Methods("GET")

	// Create new subrouter so we can wrap the JWT middleware around it
	wsAuth := ws.PathPrefix("/").Subrouter()
	wsAuth.Use(middlewares.ValidateJWT)
	wsAuth.HandleFunc("/v1/auth/check", WebsocketAuthCheckV1).Methods("GET")
}

func WebsocketInfo(w http.ResponseWriter, r *http.Request) {
	helloWorld := responses.Generic{
		Message: "kees websocket portal",
		Data: map[string]interface{}{
			"available": map[string]string{
				"mediacontroller": "/v1/mc",
				"webclient":       "/v1/wc",
			},
		},
	}

	data, err := helpers.Format(helloWorld)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

func WebsocketAuthV1(w http.ResponseWriter, r *http.Request) {
	// validate the auth header for the connecting media controller
	token := r.Header.Get("X-Kees-MC-Token")
	helpers.Debug(token)

	device, err := models.Devices.ByToken(token)
	if device == nil {
		data, err := helpers.Format(responses.Generic{
			Message: "Unauthorized Media Controller",
			Data:    map[string]interface{}{},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	deviceUpdate := DeviceUpdate{}
	err = helpers.Parse(r, &deviceUpdate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(deviceUpdate.Version) == 0 {
		data, err := helpers.Format(responses.Generic{
			Message: "Invalid Media Controller Version",
			Data:    map[string]interface{}{},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	device.Version = deviceUpdate.Version
	device.Capabilities = deviceUpdate.Capabilities
	err = device.Update()

	jwt, expiresIn, err := helpers.GenerateJWT(map[string]interface{}{
		"id":           device.ID,
		"name":         device.Name,
		"version":      device.Version,
		"controller":   device.Controller,
		"capabilities": device.Capabilities,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.Debug(jwt)

	jwtResponse := AuthResponse{
		Message: "successfully authd " + device.Name,
		Device:  *device,
		JWT: JWTResponse{
			ExpiresIn: expiresIn,
			Token:     jwt,
		},
	}

	data, err := helpers.Format(jwtResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

func WebsocketAuthCheckV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	data, err := helpers.Format(jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
	return
}
