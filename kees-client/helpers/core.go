package helpers

import (
	"kees-client/config"
)

func Configure(config *config.Config) {
	ConfigureLogging(config.Logging)
}
