package helpers

import (
	"kees-server/config"
)

// TODO: add proper error returns to ConfigureX helpers
func Configure(config *config.Config) {
	ConfigureLogging(config.Logging)
	ConfigureJWT(config.JWT)
}
