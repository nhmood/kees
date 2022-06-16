package web

import (
	"net/http"

	"github.com/Masterminds/log-go"
	"github.com/gorilla/mux"

	"kees/server/config"
	"kees/server/constants"
	"kees/server/devices"
	"kees/server/helpers"

	"kees/server/web/api"
	"kees/server/web/middlewares"
	"kees/server/web/websocket"
)

var handler http.Handler
var serverConfig config.ServerConfig

func Configure(c config.ServerConfig) {
	log.Info("Configuring kees-server API")
	serverConfig = c

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Root)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/")))
	router.PathPrefix("/static").Handler(fs)

	broker := devices.NewBroker()
	broker.Run()

	api.Configure(router, "/api", broker)
	websocket.Configure(router, "/ws", broker)

	handler = middlewares.AddBaseHeaders(router)
	handler = middlewares.AddLogging(handler)
	handler = middlewares.CORSHandler(handler)
}

func Root(w http.ResponseWriter, r *http.Request) {
	// root endpoint doesn't have its own subrouter so we can't apply a middleware
	// for this one endpoint, manually set the content-type to application/json
	w.Header().Add("Content-Type", "application/json")
	helpers.Halt(w, http.StatusOK, "Hello world from kees-server", map[string]interface{}{
		"commit": constants.GitCommit,
	})
	return
}

func Run() {
	log.Info("Starting Server on:" + serverConfig.Port)
	log.Fatal(http.ListenAndServe(":"+serverConfig.Port, handler))
}
