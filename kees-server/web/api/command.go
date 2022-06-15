package api

import (
	"net/http"

	"github.com/Masterminds/log-go"

	"kees/server/devices"
	"kees/server/helpers"
	"kees/server/models"

	"kees/server/web/responses"
)

type Command struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type MCResponse struct {
	devices.MediaControllerInfo
	Capabilities []string `json:"capabilities"`
}

type CommandResponse struct {
	Device  MCResponse `json:"device"`
	Command Command    `json:"command"`
}

//type CommandResponse struct {
//	Device   MCResponse               `json:"device"`
//	Commands []map[string]interface{} `json:"commands"`
//}
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
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	deviceID := helpers.GetStringParam(r, "device_id", helpers.URLParam)
	device, err := models.Devices.Get(deviceID)

	if device == nil {
		data, err := helpers.Format(responses.Generic{
			Message: "DeviceID: " + deviceID + " not found",
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

	if !device.Online {
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
	validOperation := device.ValidOperation(operation)

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

	mc := broker.MediaControllers[device.ID]
	if mc == nil {
		data, err := helpers.Format(responses.Generic{
			Message: "DeviceID: " + deviceID + " broker connection missing",
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