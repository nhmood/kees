package api

import (
	"net/http"

	"kees-server/devices"
	"kees-server/helpers"

	"kees-server/web/responses"
)

type MCResponse struct {
	devices.MediaControllerInfo
	Capabilities []string `json:"capabilities"`
}

func DevicesV1(w http.ResponseWriter, r *http.Request) {
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
	helpers.Debug(jwt)

	mcs := make([]*MCResponse, 0)
	for _, v := range broker.MediaControllers {
		// TODO: pull capabilities by device type from database
		mcr := &MCResponse{MediaControllerInfo: v.Info, Capabilities: []string{
			"play",
			"stop",
			"rewind",
			"fast_forward",
			"pause",
			"shuffle",
		}}

		mcs = append(mcs, mcr)
	}

	data, err := helpers.Format(mcs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(data)
	return
}
