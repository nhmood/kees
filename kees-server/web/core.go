package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"kees-server/config"
	"kees-server/constants"
	"kees-server/devices/broker"
	"kees-server/helpers"

	"kees-server/web/api"
	"kees-server/web/middlewares"
	"kees-server/web/responses"
	"kees-server/web/websocket"
)

var handler http.Handler
var serverConfig config.ServerConfig
var connBroker *broker.Broker

func Configure(c config.ServerConfig) {
	log.Println("Configuring kees-server API")
	serverConfig = c

	connBroker = broker.New()
	helpers.Dump(connBroker)
	connBroker.Run()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Root)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/")))
	router.PathPrefix("/static").Handler(fs)

	api.Configure(router, "/api", connBroker)
	websocket.Configure(router, "/ws", connBroker)

	handler = middlewares.AddBaseHeaders(router)
	handler = middlewares.AddLogging(handler)
	handler = middlewares.CORSHandler(handler)
}

func Root(w http.ResponseWriter, r *http.Request) {
	helloWorld := responses.Generic{
		Message: "Hello World from kees-server!",
		Data: map[string]interface{}{
			"commit": constants.GitCommit,
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

func Run() {
	log.Print("Starting Server on:" + serverConfig.Port)
	log.Fatal(http.ListenAndServe(":"+serverConfig.Port, handler))
}
