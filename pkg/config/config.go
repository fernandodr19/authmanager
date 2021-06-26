package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config defines the service configuration
type Config struct {
	AppName string `envconfig:"APP_NAME" default:"authenticator"`
	API     APIConfig
}

type APIConfig struct {
	Address         string        `envconfig:"SERVER_ADDRESS" default:"0.0.0.0:3000"`
	SwaggerHost     string        `envconfig:"SWAGGER_HOST" default:"0.0.0.0:3001"`
	ShutdownTimeout time.Duration `envconfig:"APP_SHUTDOWN_TIMEOUT" default:"5s"`
}

func Load() (*Config, error) {
	var config Config
	noPrefix := ""
	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
