package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	app "github.com/fernandodr19/library/pkg"
	"github.com/fernandodr19/library/pkg/config"
	"github.com/fernandodr19/library/pkg/gateway/api"
	"github.com/fernandodr19/library/pkg/gateway/api/accounts"
	"github.com/fernandodr19/library/pkg/gateway/authorizer"
	"github.com/fernandodr19/library/pkg/gateway/repositories"
	"github.com/fernandodr19/library/pkg/instrumentation/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	_ "github.com/joho/godotenv/autoload"
)

type testEnviroment struct {
	Server *httptest.Server
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
		logger.WithError(err).Panic("failed loading config")
	}

	pgContainer := testcontainers.NewLocalDockerCompose(
		[]string{"./docker-compose.yml"},
		strings.ToLower(uuid.New().String()),
	)
	pgContainer.WithCommand([]string{"up", "-d"}).Invoke()
	// pgContainer.W([]string{"database system is ready to accept connections"})
	time.Sleep(5 * time.Second)

	// Init postgres
	dbConn, err := repositories.NewConnection(cfg.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("failed setting up postgres")
	}

	// Init authorizer
	auth, err := authorizer.New(cfg.API.TokenSecret)
	if err != nil {
		logger.WithError(err).Panic("failed building authorizer")
	}

	// Build app
	app, err := app.BuildApp(dbConn, cfg, auth)
	if err != nil {
		logger.WithError(err).Panic("failed building app")
	}

	// Build API handler
	apiHandler, err := api.BuildHandler(app, cfg, auth)
	if err != nil {
		logger.WithError(err).Panic("Could not initialize api")
	}

	testEnv.Server = httptest.NewServer(apiHandler)

	return func() {
		pgContainer.Down()
		time.Sleep(1 * time.Second)
	}
}

// https://medium.com/trendyol-tech/kafka-test-containers-with-golang-b85e4b2469db

func Test_SignUp(t *testing.T) {
	target := testEnv.Server.URL + "/api/v1/accounts/signup"
	body, err := json.Marshal(
		accounts.CreateAccountRequest{
			Email:    "aaat@test.com",
			Password: "123",
		})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
	require.NoError(t, err)

	// test
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// assert
	require.Equal(t, http.StatusCreated, resp.StatusCode)

}
