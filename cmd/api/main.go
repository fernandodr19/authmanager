package main

import (
	"net/http"
	"time"

	"github.com/fernandodr19/authmanager/docs/swagger"
	_ "github.com/fernandodr19/authmanager/docs/swagger"
	app "github.com/fernandodr19/authmanager/pkg"
	"github.com/fernandodr19/authmanager/pkg/config"
	"github.com/fernandodr19/authmanager/pkg/gateway/api"
	"github.com/fernandodr19/authmanager/pkg/gateway/authorizer"
	"github.com/fernandodr19/authmanager/pkg/gateway/repositories"
	"github.com/fernandodr19/authmanager/pkg/instrumentation/logger"

	_ "github.com/joho/godotenv/autoload"
)

// @title Swagger Authorization Manager API
// @version 1.0
// @host localhost:3000
// @basePath /api/v1
// @schemes http
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @description Documentation Authorization Manager API
func main() {
	logger := logger.Default()
	logger.Infof("build info: time[%s] git_hash[%s]", BuildTime, BuildGitCommit)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed loading config")
	}

	swagger.SwaggerInfo.Host = cfg.Swagger.Host

	// Init postgres
	dbConn, err := repositories.NewConnection(cfg.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("failed setting up postgres")
	}

	// Init authorizer
	auth, err := authorizer.New(cfg.API.TokenSecret)
	if err != nil {
		logger.WithError(err).Fatal("failed building authorizer")
	}

	// Build app
	app, err := app.BuildApp(dbConn, cfg, auth)
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
		Addr:         cfg.API.Address(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Default().WithField("address", cfg.API.Address()).Info("server starting...")
	logger.Default().Fatal(server.ListenAndServe())

}
