package helpers

import (
	"kees/media-controller/config"
)

func Configure(config *config.Config) {
	ConfigureLogging(config.Logging)
}
