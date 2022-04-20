package websocket

import (
	"net/http"

	"kees-server/devices"
	"kees-server/helpers"
)

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

	mediaController := devices.NewMediaController(c, broker)
	helpers.Debug(mediaController)
	mediaController.Run()
}
