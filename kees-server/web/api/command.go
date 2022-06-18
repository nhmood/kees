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

type CommandAllList struct {
	Commands  []*models.Command `json:"commands"`
	Page      int64             `json:"page"`
	PageCount int64             `json:"page_count"`
}

type CommandByDeviceList struct {
	CommandAllList
	Device *models.Device `json:"device"`
}

type CommandResponse struct {
	Device  models.Device  `json:"device"`
	Command models.Command `json:"command"`
}

func CommandHistoryV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	page := helpers.GetPage(r)
	commands, err := models.Commands.All(page)

	if err != nil {
		data := helpers.ToInterface(err)
		helpers.Halt(w, http.StatusInternalServerError, "Failed to lookup Command List", data)
		return
	}

	resp := CommandAllList{
		Commands:  commands,
		Page:      page,
		PageCount: models.Commands.PageCount,
	}

	data := helpers.ToInterface(resp)
	helpers.Halt(w, http.StatusOK, "Command History", data)
	return
}

func CommandHistoryByDeviceV1(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(map[string]interface{})
	helpers.Debug(jwt)

	page := helpers.GetPage(r)

	deviceID := helpers.GetStringParam(r, "device_id", helpers.URLParam)
	helpers.Dump(deviceID)

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
	commands, err := models.Commands.ByDevice(deviceID, page)

	resp := CommandByDeviceList{
		CommandAllList: CommandAllList{
			Commands:  commands,
			Page:      page,
			PageCount: models.Commands.PageCount,
		},
		Device: device,
	}

	data := helpers.ToInterface(resp)
	helpers.Halt(w, http.StatusOK, "Command History", data)
	return
}

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
