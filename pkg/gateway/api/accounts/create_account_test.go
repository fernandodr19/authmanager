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

	acc "github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/gateway/api/shared"
)

const JSONContentType = "application/json"

func TestHandler_CreateAccount(t *testing.T) {
	const (
		routePattern = "/api/v1/signup"
		target       = "/api/v1/signup"
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
			Handler: createHandler(acc.ErrInvalidEmail),
			Req: CreateAccountRequest{
				Email:    "invalid",
				Password: "123",
			},
			ExpectedStatusCode:   http.StatusBadRequest,
			ExpectedErrorPayload: responses.ErrInvalidEmail,
		},
		{
			Name:    "bad req invalid email",
			Handler: createHandler(acc.ErrInvalidPassword),
			Req: CreateAccountRequest{
				Email:    "valid@gmail.com",
				Password: "",
			},
			ExpectedStatusCode:   http.StatusBadRequest,
			ExpectedErrorPayload: responses.ErrInvalidPassword,
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
			assert.NotEmpty(t, response.Header().Get(shared.XReqID))
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
		Usecase: &accounts.AccountsMockUsecase{
			CreateAccountFunc: func(in1 context.Context, in2 vos.Email, in3 vos.Password) error {
				return err
			},
		},
	}
}
