package config

type ServerConfig struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Token string `yaml:"token"`
}
