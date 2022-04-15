package config

type DatabaseConfig struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}
