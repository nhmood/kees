package api

import (
	"github.com/Masterminds/log-go"
	"net/http"

	"kees-server/devices"
	"kees-server/helpers"

	"kees-server/web/responses"
)

type MCResponse struct {
	devices.MediaControllerInfo
	Capabilities []string `json:"capabilities"`
}

type CommandResponse struct {
	Device  MCResponse `json:"device"`
	Command Command    `json:"command"`
}

type Command struct {
	ID      string `json:"id"`
	Command string `json:"command"`
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

func DeviceInfoV1(w http.ResponseWriter, r *http.Request) {
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

	// TODO: add database lookup along with broker lookup
	deviceID := helpers.GetStringParam(r, "device_id", helpers.URLParam)
	mc := broker.MediaControllers[deviceID]

	if mc == nil {
		data, err := helpers.Format(responses.Generic{
			Message: "DeviceID: " + deviceID + " not online",
			Data:    map[string]interface{}{},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	mcr := &MCResponse{MediaControllerInfo: mc.Info, Capabilities: []string{
		"play",
		"stop",
		"rewind",
		"fast_forward",
		"pause",
		"shuffle",
	}}

	data, err := helpers.Format(mcr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(data)
	return
}

//func CommandHistoryV1(w http.ResponseWriter, r *http.Request) {
//	token := r.Header.Get("X-Kees-JWT-Token")
//	helpers.Debug(token)
//
//	jwt, err := helpers.ValidateJWT(token)
//	if err != nil {
//		data, err := helpers.Format(responses.Generic{
//			Message: "Invalid JWT",
//		})
//
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusBadRequest)
//		w.Write(data)
//		return
//	}
//	helpers.Debug(jwt)
//
//	data, err := helpers.Format(responses.Generic{
//		Message: "Not implemented",
//		Data:    map[string]interface{}{},
//	})
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	w.WriteHeader(http.StatusBadRequest)
//	w.Write(data)
//	return
//}

func CommandIssueV1(w http.ResponseWriter, r *http.Request) {
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

	// TODO: add database lookup along with broker lookup
	deviceID := helpers.GetStringParam(r, "device_id", helpers.URLParam)
	mc := broker.MediaControllers[deviceID]

	if mc == nil {
		data, err := helpers.Format(responses.Generic{
			Message: "DeviceID: " + deviceID + " not online",
			Data:    map[string]interface{}{},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	operation := helpers.GetStringParam(r, "operation", helpers.URLParam)

	// TODO: pull capabilities by device type from database
	validOperation := func(string) bool {
		capabilities := []string{
			"play",
			"stop",
			"rewind",
			"fast_forward",
			"pause",
			"shuffle",
		}

		for _, c := range capabilities {
			if c == operation {
				return true
			}
		}
		return false
	}(operation)

	if !validOperation {
		data, err := helpers.Format(responses.Generic{
			Message: "Operation: " + operation + " not valid capability",
			Data:    map[string]interface{}{},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	// TODO: create command record in database
	commandID := mc.IssueCommand(operation)
	log.Info("Command:" + commandID + "/" + operation + " created for " + mc.Info.ID)
	helpers.Debug(commandID)

	resp := CommandResponse{
		Device: MCResponse{
			MediaControllerInfo: mc.Info,
			Capabilities: []string{
				"play",
				"stop",
				"rewind",
				"fast_forward",
				"pause",
				"shuffle",
			},
		},
		Command: Command{
			ID:      commandID,
			Command: operation,
		},
	}

	data, err := helpers.Format(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(data)
	return
}
