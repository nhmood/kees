package api

import (
	"github.com/gorilla/mux"
	"net/http"

	"kees/server/devices"
	"kees/server/helpers"

	"kees/server/web/middlewares"
)

var broker *devices.Broker

func Configure(router *mux.Router, path string, b *devices.Broker) {
	broker = b
	api := router.PathPrefix(path).Subrouter()
	api.Use(middlewares.AddJSONHeader)

	api.HandleFunc("/", Root).Methods("GET")
	api.HandleFunc("/v1/auth", ClientAuthV1).Methods("POST")

	// Create new subrouter so we can wrap the JWT middleware around it
	apiAuth := api.PathPrefix("/").Subrouter()
	apiAuth.Use(middlewares.ValidateJWT)
	apiAuth.HandleFunc("/v1/devices", DevicesV1).Methods("GET")
	apiAuth.HandleFunc("/v1/devices/add", DeviceAddV1).Methods("POST")
	apiAuth.HandleFunc("/v1/devices/{device_id}", DeviceInfoV1).Methods("GET")
	apiAuth.HandleFunc("/v1/devices/{device_id}/delete", DeviceDeleteV1).Methods("POST")
	apiAuth.HandleFunc("/v1/devices/{device_id}/commands/{operation}", CommandIssueV1).Methods("POST")
}

func Root(w http.ResponseWriter, r *http.Request) {
	helpers.Halt(w, http.StatusOK, "Hello World!", nil)
	return
}
