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

	usecase "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/gateway/api/shared"
)

func TestHandler_Login(t *testing.T) {
	const (
		routePattern = "/api/v1/accounts/login"
		target       = "/api/v1/accounts/login"
	)

	request := func(body []byte) *http.Request {
		return httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	}

	testTable := []struct {
		Name                 string
		Handler              Handler
		Req                  LoginRequest
		ExpectedStatusCode   int
		ExpectedErrorPayload responses.ErrorPayload
	}{
		{
			Name:    "happy path",
			Handler: loginHandler(nil),
			Req: LoginRequest{
				Email:    "valid@gmail.com",
				Password: "123",
			},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			Name:    "user not found",
			Handler: loginHandler(usecase.ErrAccountNotFound),
			Req: LoginRequest{
				Email:    "valid@gmail.com",
				Password: "123",
			},
			ExpectedStatusCode:   http.StatusNotFound,
			ExpectedErrorPayload: responses.ErrAccountNotFound,
		},
		{
			Name:    "wrong password",
			Handler: loginHandler(usecase.ErrWrongPassword),
			Req: LoginRequest{
				Email:    "valid@gmail.com",
				Password: "123",
			},
			ExpectedStatusCode:   http.StatusUnauthorized,
			ExpectedErrorPayload: responses.ErrWrongPassword,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			// prepare
			response := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc(routePattern, middleware.Handle(tt.Handler.Login)).Methods(http.MethodPost)

			body, err := json.Marshal(tt.Req)
			require.NoError(t, err)

			// test
			router.ServeHTTP(response, request(body))

			//assert
			assert.Equal(t, tt.ExpectedStatusCode, response.Code)
			assert.NotEmpty(t, response.Header().Get(shared.XReqID))
			assert.Equal(t, JSONContentType, response.Header().Get("content-type"))

			if response.Code != http.StatusOK {
				var errPayload responses.ErrorPayload
				err = json.NewDecoder(response.Body).Decode(&errPayload)
				require.NoError(t, err)
				assert.Equal(t, tt.ExpectedErrorPayload, errPayload)
			}
		})
	}
}

func loginHandler(err error) Handler {
	return Handler{
		Usecase: &usecase.AccountsMockUsecase{
			LoginFunc: func(in1 context.Context, in2 vos.Email, in3 vos.Password) (vos.Tokens, error) {
				return vos.Tokens{}, err
			},
		},
	}
}
