package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"kees-server/helpers"
	"kees-server/web/responses"
)

var upgrader = websocket.Upgrader{}

type MediaControllerInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Controller string `json:"controller"`
}

type AuthResponse struct {
	Message string              `json:"message"`
	Device  MediaControllerInfo `json:"device"`
	JWT     JWTResponse         `json:"jwt"`
}

type JWTResponse struct {
	ExpiresIn int64  `json:"expires_in"`
	Token     string `json:"token"`
}

func Configure(router *mux.Router, path string) {
	ws := router.PathPrefix(path).Subrouter()
	ws.HandleFunc("/", WebsocketInfo).Methods("GET")
	ws.HandleFunc("/v1/auth", WebsocketAuthV1).Methods("POST")
	ws.HandleFunc("/v1/auth/check", WebsocketAuthCheckV1).Methods("GET")
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

	// TODO: make this a database lookup
	if token != "cdplayer" {
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

	controllerInfo := MediaControllerInfo{}
	err := helpers.Parse(r, &controllerInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	helpers.Debug(controllerInfo)
	jwt, expiresIn, err := helpers.GenerateJWT(map[string]string{
		"name":       controllerInfo.Name,
		"version":    controllerInfo.Version,
		"controller": controllerInfo.Controller,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.Debug(jwt)

	jwtResponse := AuthResponse{
		Message: "successfully authd" + controllerInfo.Name,
		Device:  controllerInfo,
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
	token := r.Header.Get("X-Kees-JWT-Token")
	helpers.Debug(token)

	jwt, err := helpers.ValidateJWT(token)
	if err != nil {
		data, err := helpers.Format(responses.Generic{
			Message: "Invalid JWT",
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	data, err := helpers.Format(jwt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
	return
}
