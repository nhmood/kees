package api

import (
	"github.com/gorilla/mux"
	"net/http"

	"kees-server/helpers"
	"kees-server/web/responses"
)

func Configure(router *mux.Router, path string) {
	router.HandleFunc("/", Root).Methods("GET")
}

func Root(w http.ResponseWriter, r *http.Request) {

	helloWorld := responses.Generic{
		Message: "Hello World!",
	}

	data, err := helpers.Format(helloWorld)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
