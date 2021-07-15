package repositories

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	"github.com/fernandodr19/authmanager/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needed to describe db driver
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/jackc/pgx/v4"
)

// NewConnection sets up a new connection with migrations
func NewConnection(cfg config.Postgres) (*pgx.Conn, error) {
	// Maybe use connection pool later on..
	conn, err := pgx.Connect(context.Background(), cfg.URL())
	if err != nil {
		return nil, err
	}

	err = runMigrations(cfg.URL())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//go:embed migrations
var migrations embed.FS

func getMigrationHandler(dbUrl string) (*migrate.Migrate, error) {
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("httpfs", source, dbUrl)
}

func runMigrations(dbUrl string) error {
	h, err := getMigrationHandler(dbUrl)
	if err != nil {
		return err
	}

	if err := h.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	srcErr, dbErr := h.Close()
	if srcErr != nil {
		return fmt.Errorf("failed to close DB source: %w", err)
	}
	if dbErr != nil {
		return fmt.Errorf("failed to close migrations repositories connection: %w", err)
	}

	return nil
}
