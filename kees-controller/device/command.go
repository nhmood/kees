package device

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/log-go"

	"kees/controller/helpers"
	"kees/controller/messages"
)

func (c *Controller) Command(payload *messages.WebSocket) *messages.WebSocket {
	helpers.Debug(payload)
	command := payload.Data["operation"].(string)
	operation, ok := c.Device.Capabilities[command]
	if !ok {
		log.Error("Unknown command ", command)
		err := formatMessage("error", "Unknown command", &payload.Data)
		return &err
	}

	runCommand(operation)

	// TODO: acknowledge command status
	return nil
}

func runCommand(command string) bool {
	log.Info("Executing ", command)

	split := strings.Split(command, " ")
	name := split[0]
	args := split[1:]

	cmd := exec.Command(name, args...)
	output, err := cmd.Output()

	if err != nil {
		helpers.Dump(err)
		return false
	}

	log.Info(fmt.Sprintf("%s -> %s", command, output))
	return true
}
