package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
	Server   ServerConfig   `yaml:"server"`
}

type LoggingConfig struct {
	Debug     bool `yaml:"debug"`
	SpewDepth int  `yaml:"spew_depth"`
}

type DatabaseConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

func ReadConfig(filename string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open %s", filename)
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Panic(err)
		return config, err
	}

	return config, nil
}
