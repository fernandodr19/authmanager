package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config defines the service configuration
type Config struct {
	AppName string `envconfig:"APP_NAME" default:"authmanager"`
	API
	Postgres
}

// API defines api configuration
type API struct {
	Address         string        `envconfig:"SERVER_ADDRESS" default:"0.0.0.0:3000"`
	SwaggerHost     string        `envconfig:"SWAGGER_HOST" default:"0.0.0.0:3001"`
	ShutdownTimeout time.Duration `envconfig:"APP_SHUTDOWN_TIMEOUT" default:"5s"`
	TokenSecret     string        `envconfig:"TOKEN_SECRET" default:"My Secret"`
}

// Postgres defines postgres configuration
type Postgres struct {
	User     string `envconfig:"DATABASE_USER" default:"postgres"`
	Password string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	Host     string `envconfig:"DATABASE_HOST_DIRECT" default:"localhost"`
	Port     string `envconfig:"DATABASE_PORT_DIRECT" default:"5432"`
	DBName   string `envconfig:"DATABASE_NAME" default:"dev"`
	SSLMode  string `envconfig:"DATABASE_SSLMODE" default:"disable"`
}

// URL builds postgres URL
func (p Postgres) URL() string {
	// example: "postgres://username:password@localhost:5432/db_name"
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.DBName,
		p.SSLMode,
	)
}

// Load loads config
func Load() (*Config, error) {
	var config Config
	noPrefix := ""
	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
