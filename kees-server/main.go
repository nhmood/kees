package main

import (
	"fmt"
	"log"
	"os"

	"kees-server/config"
	"kees-server/helpers"
	"kees-server/models"

	"kees-server/web"
)

func main() {
	info()
	configPath := os.Getenv("KEES_CONFIG")
	Configuration, err := config.ReadConfig(configPath)
	if err != nil {
		log.Print("Failed to read config yaml from KEES_CONFIG -> ", configPath)
		log.Panic(err)
		os.Exit(1)
	}

	helpers.Configure(Configuration)

	err = models.Configure(Configuration.Database)
	helpers.Debug(err)

	web.Configure(Configuration.Server)
	web.Run()
}

func info() {
	log.Print("kees v0.0.1")
}

func help() {
	info()

	fmt.Println(``)
}
