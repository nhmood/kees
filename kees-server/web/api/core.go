package api

import (
	"github.com/gorilla/mux"
	"net/http"

	"kees-server/devices"
	"kees-server/helpers"

	"kees-server/web/middlewares"
	"kees-server/web/responses"
)

var broker *devices.Broker

func Configure(router *mux.Router, path string, b *devices.Broker) {
	broker = b
	api := router.PathPrefix(path).Subrouter()
	api.HandleFunc("/", Root).Methods("GET")
	api.HandleFunc("/v1/auth", ClientAuthV1).Methods("POST")

	api.HandleFunc("/v1/devices", DevicesV1).Methods("GET")
	api.HandleFunc("/v1/devices/{device_id}", DeviceInfoV1).Methods("GET")

	api.Use(middlewares.AddJSONHeader)
}

func Root(w http.ResponseWriter, r *http.Request) {

	helloWorld := responses.Generic{
		Message: "Hello World!",
		Data:    map[string]interface{}{},
	}

	data, err := helpers.Format(helloWorld)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
