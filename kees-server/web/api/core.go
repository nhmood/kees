package api

import (
	"github.com/gorilla/mux"
	"net/http"

	"kees-server/devices/broker"
	"kees-server/helpers"
	"kees-server/web/middlewares"
	"kees-server/web/responses"
)

var connBroker *broker.Broker

func Configure(router *mux.Router, path string, broker *broker.Broker) {
	connBroker = broker

	api := router.PathPrefix(path).Subrouter()
	api.HandleFunc("/", Root).Methods("GET")

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
