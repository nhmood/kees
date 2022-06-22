package device

import (
	"fmt"
	"os/exec"

	"github.com/Masterminds/log-go"

	"kees/controller/helpers"
	"kees/controller/messages"
)

type CommandHandler func(*Controller) *messages.WebSocket

var CommandHandlers = map[string]CommandHandler{
	"play":         (*Controller).Play,
	"stop":         (*Controller).Stop,
	"rewind":       (*Controller).Rewind,
	"fast_forward": (*Controller).FastForward,
	"pause":        (*Controller).Pause,
	"shuffle":      (*Controller).Shuffle,
}

func (c *Controller) Command(payload *messages.WebSocket) *messages.WebSocket {
	helpers.Debug(payload)
	command := payload.Data["command"].(string)
	commandFunc := CommandHandlers[command]
	if commandFunc == nil {
		log.Error("Unknown command ", command)
		err := formatMessage("error", "Unknown command", &payload.Data)
		return &err
	}

	commandFunc(c)

	// TODO: acknowledge command status
	return nil
}

func runCommand(command string) bool {
	log.Info("Executing bin/irsend ", command)
	cmd := exec.Command("bin/irsend", command)
	output, err := cmd.Output()

	if err != nil {
		helpers.Dump(err)
		return false
	}

	log.Info(fmt.Sprintf("irsend-> %s", output))
	return true
}

func (c *Controller) Play() *messages.WebSocket {
	runCommand("play")
	return nil
}

func (c *Controller) Stop() *messages.WebSocket {
	runCommand("stop")
	return nil
}

func (c *Controller) Rewind() *messages.WebSocket {
	runCommand("rewind")
	return nil
}

func (c *Controller) FastForward() *messages.WebSocket {
	runCommand("fast_forward")
	return nil
}

func (c *Controller) Pause() *messages.WebSocket {
	runCommand("pause")
	return nil
}

func (c *Controller) Shuffle() *messages.WebSocket {
	runCommand("shuffle")
	return nil
}
