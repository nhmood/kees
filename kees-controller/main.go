package main

import (
	"github.com/Masterminds/log-go"
	"os"

	"kees/controller/config"
	"kees/controller/device"
	"kees/controller/helpers"
)

func main() {
	info()
	configPath := os.Getenv("KEES_CONFIG")
	Configuration, err := config.ReadConfig(configPath)
	if err != nil {
		log.Error("Failed to read config yaml from KEES_CONFIG:", configPath)
		os.Exit(1)
	}

	helpers.Configure(Configuration)

	device := device.NewController(Configuration)
	device.Run()
}

func info() {
	log.Info("kees controller v0.0.1")
}

func help() {
	info()

	log.Info(``)
}
