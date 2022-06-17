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
	Device  models.Device  `json:"device"`
	Command models.Command `json:"command"`
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

	if err != nil {
		data := helpers.ToInterface(err)
		helpers.Halt(w, http.StatusInternalServerError, "Device Lookup Failed", data)
		return
	}

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

	command := &models.Command{
		Operation: operation,
		Status:    "new",
		Metadata:  "",
		Client:    r.UserAgent(),
		DeviceID:  device.ID,
	}
	helpers.Dump(command)

	command, err = models.Commands.Insert(*command)
	helpers.Dump(command)

	// TODO: create command record in database
	mc.IssueCommand(command)
	log.Info("Command:" + command.ID + "/" + command.Operation + " created for " + command.DeviceID)
	helpers.Debug(command.ID)

	resp := CommandResponse{
		Device:  *device,
		Command: *command,
	}

	data := helpers.ToInterface(resp)
	helpers.Halt(w, 200, "Command: "+command.Operation+" initiated", data)
	return
}
