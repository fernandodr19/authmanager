package tests

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	app "github.com/fernandodr19/authmanager/pkg"
	"github.com/fernandodr19/authmanager/pkg/config"
	"github.com/fernandodr19/authmanager/pkg/gateway/api"
	"github.com/fernandodr19/authmanager/pkg/gateway/authorizer"
	"github.com/fernandodr19/authmanager/pkg/gateway/repositories"
	"github.com/fernandodr19/authmanager/pkg/instrumentation/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/testcontainers/testcontainers-go"

	_ "github.com/joho/godotenv/autoload"
)

type testEnviroment struct {
	Server *httptest.Server
	Conn   *pgx.Conn
	App    *app.App
}

var testEnv testEnviroment

func TestMain(m *testing.M) {
	teardown := setup()
	exitCode := m.Run()
	teardown()
	os.Exit(exitCode)
}

func setup() func() {
	logger := logger.Default()
	logger.Info("setting up integration tests env")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("failed loading config")
	}

	err = setupDocker()
	if err != nil {
		logger.WithError(err).Fatal("failed setting up docker")
	}

	// Setup postgres
	cfg.Postgres.DBName = "test"
	cfg.Postgres.Port = "5433"
	dbConn, err := setupPostgres(cfg.Postgres)
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

	testEnv.Conn = dbConn
	testEnv.App = app
	testEnv.Server = httptest.NewServer(apiHandler)

	return func() {
		TruncatePostgresTables()
	}
}

// https://medium.com/trendyol-tech/kafka-test-containers-with-golang-b85e4b2469db

func setupDocker() error {
	running, err := isDockerRunning([]string{
		"pg-test",
	})
	if err != nil {
		return err
	}

	if running {
		logger.Default().Infoln("necessary containers already running...")
		return nil
	}

	compose := testcontainers.NewLocalDockerCompose(
		[]string{"./docker-compose.yml"},
		strings.ToLower(uuid.New().String()),
	)
	execErr := compose.WithCommand([]string{"up", "-d"}).Invoke()
	if execErr.Error != nil {
		return execErr.Error
	}
	return nil
}

func isDockerRunning(expectedImages []string) (bool, error) {
	stdout, err := exec.Command("docker", "ps").Output()
	if err != nil {
		return false, err
	}

	ps := string(stdout)
	if err != nil {
		return false, err
	}

	running := true
	for _, image := range expectedImages {
		if !strings.Contains(ps, image) {
			running = false
			break
		}
	}
	return running, nil
}

func setupPostgres(cfg config.Postgres) (*pgx.Conn, error) {
	done := make(chan bool, 1)
	var dbConn *pgx.Conn
	var err error

	// tries to connect within 5 seconds timeout
	go func() {
		for {
			dbConn, err = repositories.NewConnection(cfg)
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			break
		}
		close(done)
	}()

	select {
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("timed out trying to set up postgres: %w", err)
	case <-done:
	}

	return dbConn, nil
}

func TruncatePostgresTables() {
	testEnv.Conn.Exec(context.Background(),
		`TRUNCATE TABLE 
			accounts
		CASCADE`,
	)
}
