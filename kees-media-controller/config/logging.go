package config

type LoggingConfig struct {
	Debug     bool `yaml:"debug"`
	SpewDepth int  `yaml:"spew_depth"`
}
