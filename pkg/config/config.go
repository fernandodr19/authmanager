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
	Port            string        `envconfig:"PORT" default:"3000"`
	SwaggerHost     string        `envconfig:"SWAGGER_HOST" default:"0.0.0.0:3001"`
	ShutdownTimeout time.Duration `envconfig:"APP_SHUTDOWN_TIMEOUT" default:"5s"`
	TokenSecret     string        `envconfig:"TOKEN_SECRET" default:"My Secret"`
}

func (a API) Address() string {
	return fmt.Sprintf(":%s", a.Port)
}

// Postgres defines postgres configuration
type Postgres struct {
	User     string `envconfig:"DATABASE_USER" default:"postgres"`
	Password string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	Host     string `envconfig:"DATABASE_HOST" default:"localhost"`
	Port     string `envconfig:"DATABASE_PORT" default:"5432"`
	DBName   string `envconfig:"DATABASE_NAME" default:"dev"`
	SSLMode  string `envconfig:"DATABASE_SSLMODE" default:"sslmode=disable"`
}

// URL builds postgres URL
func (p Postgres) URL() string {
	// example: "postgres://username:password@localhost:5432/db_name"
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.DBName,
	)

	// built separately since some systems don't support it
	if p.SSLMode != "" {
		url = fmt.Sprintf("%s?%s", url, p.SSLMode)
	}

	return url
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
