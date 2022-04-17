package websocket

import (
	"net/http"

	"kees-server/helpers"
	"kees-server/web/responses"
)

func ClientV1(w http.ResponseWriter, r *http.Request) {
	// validate the auth header for the connecting media controller
	token := r.Header.Get("X-Kees-MC-Token")
	helpers.Dump(token)

	// TODO: make this a database lookup
	if token == "cdplayer" {
		data, err := helpers.Format(responses.Generic{
			Message: "Unauthorized Media Controller",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	// if valid, establish websocket connection
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		helpers.Dump("Failed to upgrade ws")
		helpers.Dump(err)
		return
	}

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			helpers.Dump(err)
			break
		}

		helpers.Dump(message)
		helpers.Dump(mt)
	}

}
