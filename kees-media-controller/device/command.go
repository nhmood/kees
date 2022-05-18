package device

import (
	"fmt"
	"os/exec"

	"github.com/Masterminds/log-go"

	"kees/media-controller/helpers"
	"kees/media-controller/messages"
)

type CommandHandler func(*MediaController) *messages.WebSocket

var CommandHandlers = map[string]CommandHandler{
	"play":         (*MediaController).Play,
	"stop":         (*MediaController).Stop,
	"rewind":       (*MediaController).Rewind,
	"fast_forward": (*MediaController).FastForward,
	"pause":        (*MediaController).Pause,
	"shuffle":      (*MediaController).Shuffle,
}

func (c *MediaController) Command(payload *messages.WebSocket) *messages.WebSocket {
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

func (c *MediaController) Play() *messages.WebSocket {
	runCommand("play")
	return nil
}

func (c *MediaController) Stop() *messages.WebSocket {
	runCommand("stop")
	return nil
}

func (c *MediaController) Rewind() *messages.WebSocket {
	runCommand("rewind")
	return nil
}

func (c *MediaController) FastForward() *messages.WebSocket {
	runCommand("fast_forward")
	return nil
}

func (c *MediaController) Pause() *messages.WebSocket {
	runCommand("pause")
	return nil
}

func (c *MediaController) Shuffle() *messages.WebSocket {
	runCommand("shuffle")
	return nil
}
