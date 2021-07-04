package main

import (
	"context"
	"net/http"
	"time"

	_ "github.com/fernandodr19/library/docs/swagger"
	app "github.com/fernandodr19/library/pkg"
	"github.com/fernandodr19/library/pkg/config"
	"github.com/fernandodr19/library/pkg/gateway/api"
	"github.com/fernandodr19/library/pkg/gateway/authorizer"
	"github.com/fernandodr19/library/pkg/gateway/repositories"
	"github.com/fernandodr19/library/pkg/instrumentation/logger"
	"github.com/jackc/pgx/v4"

	_ "github.com/joho/godotenv/autoload"
)

// @title Swagger library API
// @version 1.0
// @host localhost:3000
// @basePath /api/v1
// @schemes http
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @description Documentation Library API
func main() {
	logger := logger.Default()
	logger.Infof("build info: time[%s] git_hash[%s]", BuildTime, BuildGitCommit)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed loading config")
	}

	// Init postgres
	pgConn, err := setupPostgres(cfg.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("failed setting up postgres")
	}

	auth, err := authorizer.New(cfg.API.TokenSecret)
	if err != nil {
		logger.WithError(err).Fatal("failed building authorizer")
	}
	// // fmt.Println(auth.CreateToken("my-user", 3000*time.Second))

	// Build app
	app, err := app.BuildApp(pgConn, cfg, auth)
	if err != nil {
		logger.WithError(err).Fatal("failed building app")
	}

	// Build API handler
	apiHandler, err := api.BuildHandler(app, cfg, auth)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize api")
	}

	serveApp(apiHandler, cfg)
}

func serveApp(apiHandler http.Handler, cfg *config.Config) {
	server := &http.Server{
		Handler:      apiHandler,
		Addr:         cfg.API.Address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Default().WithField("address", cfg.API.Address).Info("server starting...")
	logger.Default().Fatal(server.ListenAndServe())
}

func setupPostgres(cfg config.Postgres) (*pgx.Conn, error) {
	// Maybe use connection pool later on..
	conn, err := pgx.Connect(context.Background(), cfg.URL())
	if err != nil {
		return nil, err
	}

	err = repositories.RunMigrations(cfg.URL())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
