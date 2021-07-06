package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/fernandodr19/library/docs/swagger"
	app "github.com/fernandodr19/library/pkg"
	"github.com/fernandodr19/library/pkg/config"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api"
	"github.com/fernandodr19/library/pkg/gateway/authorizer"
	"github.com/fernandodr19/library/pkg/gateway/repositories"
	"github.com/fernandodr19/library/pkg/instrumentation/logger"

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
	pgConn, err := repositories.NewConnection(cfg.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("failed setting up postgres")
	}

	// Init authorizer
	auth, err := authorizer.New(cfg.API.TokenSecret)
	if err != nil {
		logger.WithError(err).Fatal("failed building authorizer")
	}
	t, _ := auth.CreateTokens(accounts.Account{ID: "3d5a5c6a-d589-4e7f-9269-f5346fc40549"}, 1*time.Second, 3000*time.Second)
	fmt.Println(t.AccessToken)

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
