package repositories

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // needed to describe db driver
	"github.com/golang-migrate/migrate/v4/source/httpfs"
)

//go:embed migrations
var migrations embed.FS

func getMigrationHandler(dbUrl string) (*migrate.Migrate, error) {
	source, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		return nil, err
	}
	return migrate.NewWithSourceInstance("httpfs", source, dbUrl)
}

// RunMigrations runs postgres migrations
func RunMigrations(dbUrl string) error {
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
