package main

import (
	"github.com/Masterminds/log-go"
	"os"

	"kees-client/config"
	"kees-client/helpers"
)

func main() {
	info()
	configPath := os.Getenv("KEES_CONFIG")
	Configuration, err := config.ReadConfig(configPath)
	if err != nil {
		log.Error("Failed to read config yaml from KEES_CONFIG -> ", configPath)
		log.Panic(err)
		os.Exit(1)
	}

	helpers.Configure(Configuration)
}

func info() {
	log.Info("kees client v0.0.1")
}

func help() {
	info()

	log.Info(``)
}
