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

	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/fernandodr19/library/pkg/gateway/api/shared"
)

const JSONContentType = "application/json"

func TestHandler_DoSomething(t *testing.T) {
	const (
		routePattern = "/api/v1/do-something"
		target       = "/api/v1/do-something"
	)

	request := func(body []byte) *http.Request {
		return httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	}

	t.Run("should return 200", func(t *testing.T) {
		// prepare
		handler := createHandler(nil)
		response := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(routePattern, middleware.Handle(handler.CreateAccount)).Methods(http.MethodPost)

		body, err := json.Marshal(CreateAccountRequest{})
		require.NoError(t, err)

		// test
		router.ServeHTTP(response, request(body))

		//assert
		assert.Equal(t, http.StatusOK, response.Code)
		assert.NotEmpty(t, response.Header().Get(shared.XReqID))
		assert.Equal(t, JSONContentType, response.Header().Get("content-type"))
	})

	t.Run("should return 400", func(t *testing.T) {
		// prepare
		handler := createHandler(accounts.ErrNotImplemented)
		response := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(routePattern, middleware.Handle(handler.CreateAccount)).Methods(http.MethodPost)

		// test
		router.ServeHTTP(response, request(nil))

		//assert
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.NotEmpty(t, response.Header().Get(shared.XReqID))
		assert.Equal(t, JSONContentType, response.Header().Get("content-type"))
	})
}

func createHandler(err error) Handler {
	return Handler{
		Usecase: &accounts.AccountsMockUsecase{
			CreateAccountFunc: func(in1 context.Context, in2 vos.Email, in3 vos.Password) (accounts.Tokens, error) {
				return accounts.Tokens{}, err
			},
		},
	}
}
