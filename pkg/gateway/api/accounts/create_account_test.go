package accounts

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fernandodr19/authmanager/pkg/domain/entities/accounts"
	usecase "github.com/fernandodr19/authmanager/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/fernandodr19/authmanager/pkg/gateway/api/middleware"
	"github.com/fernandodr19/authmanager/pkg/gateway/api/responses"
)

const JSONContentType = "application/json"

func TestHandler_CreateAccount(t *testing.T) {
	const (
		routePattern = "/api/v1/accounts"
		target       = "/api/v1/accounts"
	)

	request := func(body []byte) *http.Request {
		return httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	}

	testTable := []struct {
		Name                 string
		Handler              Handler
		Req                  CreateAccountRequest
		ExpectedStatusCode   int
		ExpectedErrorPayload responses.ErrorPayload
	}{
		{
			Name:    "sign up happy path",
			Handler: createHandler(nil),
			Req: CreateAccountRequest{
				Email:    "valid@gmail.com",
				Password: "123",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			Name:    "bad req invalid email",
			Handler: createHandler(accounts.ErrInvalidEmail),
			Req: CreateAccountRequest{
				Email:    "invalid",
				Password: "123",
			},
			ExpectedStatusCode:   http.StatusUnprocessableEntity,
			ExpectedErrorPayload: responses.ErrInvalidEmail,
		},
		{
			Name:    "bad req invalid password",
			Handler: createHandler(accounts.ErrInvalidPassword),
			Req: CreateAccountRequest{
				Email:    "valid@gmail.com",
				Password: "",
			},
			ExpectedStatusCode:   http.StatusUnprocessableEntity,
			ExpectedErrorPayload: responses.ErrInvalidPassword,
		},
		{
			Name:    "conflicted email",
			Handler: createHandler(usecase.ErrEmailAlreadyRegistered),
			Req: CreateAccountRequest{
				Email:    "valid@gmail.com",
				Password: "123",
			},
			ExpectedStatusCode:   http.StatusConflict,
			ExpectedErrorPayload: responses.ErrEmailAlreadyRegistered,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			// prepare
			response := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc(routePattern, middleware.Handle(tt.Handler.CreateAccount)).Methods(http.MethodPost)

			body, err := json.Marshal(tt.Req)
			require.NoError(t, err)

			// test
			router.ServeHTTP(response, request(body))

			//assert
			assert.Equal(t, tt.ExpectedStatusCode, response.Code)
			assert.Equal(t, JSONContentType, response.Header().Get("content-type"))

			if response.Code != http.StatusCreated {
				var errPayload responses.ErrorPayload
				err = json.NewDecoder(response.Body).Decode(&errPayload)
				require.NoError(t, err)
				assert.Equal(t, tt.ExpectedErrorPayload, errPayload)
			}
		})
	}
}

func createHandler(err error) Handler {
	return Handler{
		Usecase: &usecase.AccountsMockUsecase{
			CreateAccountFunc: func(in1 context.Context, in2 vos.Email, in3 vos.Password) (vos.AccID, error) {
				return "", err
			},
		},
	}
}
