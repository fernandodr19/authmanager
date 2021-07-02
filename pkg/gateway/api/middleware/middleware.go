package middleware

import (
	"net/http"
	"strings"

	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/gateway/api/shared"
	"github.com/fernandodr19/library/pkg/instrumentation"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Authorizer interface {
	AuthorizeRequest(h http.Handler) http.Handler
}

// Cors applies cors rules to router
func Cors(r *mux.Router) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	return handlers.CORS(originsOk, headersOk, methodsOk)(r)
}

// Removes the trailing slash from request, except if it is the root url.
// If the url is https://www.google.com/api or https://www.google.com/api/
// both will match.
// This was done as gorilla mux default method for this doesn't support POST requests: https://github.com/gorilla/mux/issues/30
func TrimSlashSuffix(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path != "/" {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	}

	next.ServeHTTP(w, r)
}

// Handle middleware function to treat rest responses.
func Handle(handler func(r *http.Request) responses.Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := instrumentation.Logger()

		response := handler(r)
		if response.Error != nil {
			err := response.Error
			logger.Error(err)
		}

		// Setting headers
		for key, value := range response.Headers() {
			w.Header().Set(key, value)
		}

		if w.Header().Get(shared.XReqID) == "" {
			w.Header().Set(shared.XReqID, uuid.NewString())
		}

		err := responses.SendJSON(w, response.Payload, response.Status)
		if err != nil {
			logger.Error(err)
		}
	}
}
