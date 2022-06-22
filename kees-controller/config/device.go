package config

type DeviceConfig struct {
	Token      string `yaml:"token"`
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Controller string `yaml:"controller"`
}
