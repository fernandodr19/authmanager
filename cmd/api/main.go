package main

import (
	"net/http"
	"time"

	_ "github.com/fernandodr19/authenticator/docs/swagger"
	authenticator "github.com/fernandodr19/authenticator/pkg"
	"github.com/fernandodr19/authenticator/pkg/config"
	"github.com/fernandodr19/authenticator/pkg/gateway/api"
	"github.com/fernandodr19/authenticator/pkg/instrumentation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// @title Swagger Authenticator API
// @version 1.0
// @description Documentation Auth API
func main() {
	logger := logrus.New()
	logger.Infof("Build info: time[%s] git_hash[%s]", BuildTime, BuildGitCommit)
	instrumentation.Register(&instrumentation.Instrumentation{Logger: logrus.NewEntry(logger)})

	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed loading config")
	}

	// init postgres
	app, err := authenticator.BuildApp(&pgxpool.Pool{}, cfg)
	if err != nil {
		logger.WithError(err).Fatal("failed building app")
	}

	apiHandler, err := api.BuildHandler(app, cfg)
	if err != nil {
		logger.WithError(err).Fatal("Could not initalize api")
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

	instrumentation.Logger().WithField("address", cfg.API.Address).Info("server starting...")
	instrumentation.Logger().Fatal(server.ListenAndServe())
}
