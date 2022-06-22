package helpers

import (
	"kees/controller/config"
)

func Configure(config *config.Config) {
	ConfigureLogging(config.Logging)
}
