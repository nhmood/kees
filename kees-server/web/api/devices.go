package api

import (
	"net/http"

	"kees/server/helpers"
	"kees/server/messages"
	"kees/server/models"
)

type DeviceAddPayloadV1 struct {
	Name       string `json:"name"`
	Controller string `json:"controller"`
}

type DeviceResponseV1 struct {
	Device models.Device     `json:"device"`
	Auth   map[string]string `json:"auth"`
}

func DeviceAddV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	payload := &DeviceAddPayloadV1{}
	err := helpers.Parse(r, payload)

	// TODO: should the name be provided to us (created by admin, static for device)?
	device := &models.Device{
		Name:       payload.Name,
		Controller: payload.Controller,
	}

	if len(device.Name) == 0 {
		helpers.Halt(w, http.StatusBadRequest, "Invalid Name provided for Media Controller", nil)
		return
	}

	// TODO: only allow registered controllers
	if len(device.Controller) == 0 {
		helpers.Halt(w, http.StatusBadRequest, "Invalid Name provided for Media Controller", nil)
		return
	}

	device, err = models.Devices.Insert(*device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := DeviceResponseV1{
		Device: *device,
		Auth: map[string]string{
			"token": device.Token,
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

func DeviceDeleteV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	deviceID := helpers.GetStringParam(r, "device_id", helpers.URLParam)
	device, err := models.Devices.Get(deviceID)
	// TODO: handle no result error here
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if device == nil {
		helpers.Halt(w, http.StatusBadRequest, "DeviceID: "+deviceID+" not found", nil)
		return
	}

	// TODO: pass message to broker to disconnect device
	device.Delete()
	mc := broker.MediaControllers[deviceID]
	if mc != nil {
		deleteMessage := messages.WebSocket{
			State:   "error",
			Message: "Device Deleted",
			Data:    map[string]interface{}{},
		}
		mc.Disconnect(deleteMessage)
	}

	helpers.Halt(w, http.StatusOK, "DeviceID: "+deviceID+" successfully deleted", nil)
	return
}

func DevicesV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	page := helpers.GetPage(r)
	devices, err := models.Devices.All(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := helpers.Format(devices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(data)
	return
}

func DeviceInfoV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	// TODO: add database lookup along with broker lookup
	deviceID := helpers.GetStringParam(r, "device_id", helpers.URLParam)
	device, err := models.Devices.Get(deviceID)

	if device == nil {
		helpers.Halt(w, http.StatusBadRequest, "DeviceID: "+deviceID+" not found", nil)
		return
	}

	data, err := helpers.Format(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(data)
	return
}
