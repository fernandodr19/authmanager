package accounts

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
)

const JSONContentType = "application/json"

func TestHandler_Create(t *testing.T) {
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
		router.HandleFunc(routePattern, middleware.Handle(handler.DoSomething)).Methods(http.MethodPost)

		// test
		router.ServeHTTP(response, request(nil))

		//assert
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, JSONContentType, response.Header().Get("content-type"))
	})

	t.Run("should return 501", func(t *testing.T) {
		// prepare
		handler := createHandler(accounts.ErrNotImplemented)
		response := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc(routePattern, middleware.Handle(handler.DoSomething)).Methods(http.MethodPost)

		// test
		router.ServeHTTP(response, request(nil))

		//assert
		assert.Equal(t, http.StatusNotImplemented, response.Code)
		assert.Equal(t, JSONContentType, response.Header().Get("content-type"))
	})
}

func createHandler(err error) Handler {
	return Handler{
		Usecase: &accounts.AccountsMockUsecase{
			DoSomethingFunc: func(in1 context.Context) error {
				return err
			},
		},
	}
}
