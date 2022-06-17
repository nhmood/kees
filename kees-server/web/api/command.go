package api

import (
	"net/http"

	"github.com/Masterminds/log-go"

	"kees/server/helpers"
	"kees/server/models"
)

type Command struct {
	ID      string `json:"id"`
	Command string `json:"command"`
}

type CommandResponse struct {
	Device  models.Device `json:"device"`
	Command Command       `json:"command"`
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
		helpers.Halt(w, http.StatusBadRequest, "DeviceID: "+deviceID+" not found", nil)
		return
	}

	helpers.Dump(device)

	if !device.Online {
		helpers.Halt(w, http.StatusBadRequest, "DeviceID: "+deviceID+" not online", nil)
		return
	}

	operation := helpers.GetStringParam(r, "operation", helpers.URLParam)
	validOperation := device.ValidOperation(operation)

	if !validOperation {
		helpers.Halt(w, http.StatusBadRequest, "Operation: "+operation+" not valid capability", nil)
		return
	}

	mc := broker.MediaControllers[device.ID]
	if mc == nil {
		helpers.Halt(w, http.StatusBadRequest, "DeviceID: "+deviceID+" broker connection missing", nil)
		return
	}

	// TODO: create command record in database
	commandID := mc.IssueCommand(operation)
	log.Info("Command:" + commandID + "/" + operation + " created for " + mc.Device.ID)
	helpers.Debug(commandID)

	resp := CommandResponse{
		Device: *device,
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
