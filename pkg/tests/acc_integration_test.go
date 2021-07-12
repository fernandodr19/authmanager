package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/fernandodr19/library/pkg/gateway/api/accounts"
	"github.com/stretchr/testify/require"
)

func Test_SignUp(t *testing.T) {
	defer TruncatePostgresTables()
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

func Test_Login(t *testing.T) {
	defer TruncatePostgresTables()

	ctx := context.Background()
	err := testEnv.App.Accounts.CreateAccount(ctx, "bbb@test.com", "32111")
	require.NoError(t, err)

	target := testEnv.Server.URL + "/api/v1/accounts/login"
	body, err := json.Marshal(
		accounts.LoginRequest{
			Email:    "bbb@test.com",
			Password: "32111",
		})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
	require.NoError(t, err)

	// test
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// assert
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
