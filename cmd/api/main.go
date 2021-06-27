package main

import (
	"net/http"
	"time"

	_ "github.com/fernandodr19/library/docs/swagger"
	library "github.com/fernandodr19/library/pkg"
	"github.com/fernandodr19/library/pkg/config"
	"github.com/fernandodr19/library/pkg/gateway/api"
	"github.com/fernandodr19/library/pkg/instrumentation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

// @title Swagger library API
// @version 1.0
// @description Documentation Library API
func main() {
	logger := logrus.New()
	logger.Infof("Build info: time[%s] git_hash[%s]", BuildTime, BuildGitCommit)
	instrumentation.Register(&instrumentation.Instrumentation{Logger: logrus.NewEntry(logger)})

	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed loading config")
	}

	// init postgres
	app, err := library.BuildApp(&pgxpool.Pool{}, cfg)
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
