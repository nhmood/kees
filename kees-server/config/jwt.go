package config

type JWTConfig struct {
	Issuer     string `yaml:"issuer"`
	SigningKey string `yaml:"signing_key"`
}
