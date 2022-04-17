package websocket

import (
	"github.com/gorilla/websocket"
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

	mc := devices.MC.New(c)

	mediaController := devices.mediaController.New(c)
	helpers.Dump(mediaController)

	mediaController.Active.Add(1)
	go mediaController.readHandler()
	go mediaController.writeHandler()

	mediaController.Active.Wait()
	helpers.Dump("Closing")
}
