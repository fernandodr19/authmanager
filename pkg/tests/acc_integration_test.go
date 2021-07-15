package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/fernandodr19/authmanager/pkg/gateway/api/accounts"
	"github.com/stretchr/testify/assert"
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
	_, err := testEnv.App.Accounts.CreateAccount(ctx, "bbb@test.com", "32111")
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

	var respBody accounts.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)

	assert.NotEmpty(t, respBody.AccessToken)
	assert.NotEmpty(t, respBody.RefreshToken)
}

func Test_GetAccDetails(t *testing.T) {
	defer TruncatePostgresTables()

	ctx := context.Background()

	email := vos.Email("ccc@test.com")
	_, err := testEnv.App.Accounts.CreateAccount(ctx, email, "32111")
	require.NoError(t, err)
	accID, tokens, err := testEnv.App.Accounts.Login(ctx, email, "32111")
	require.NoError(t, err)

	target := testEnv.Server.URL + "/api/v1/accounts/" + accID.String()

	req, err := http.NewRequest(http.MethodGet, target, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	require.NoError(t, err)

	// test
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// assert
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var respBody accounts.GetAccountResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)

	assert.Equal(t, accID, respBody.AccountID)
	assert.Equal(t, email, respBody.Email)
	assert.NotEmpty(t, respBody.CratedAt)
}
