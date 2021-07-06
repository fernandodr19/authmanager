package middleware

import (
	"net/http"
	"strings"

	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/gateway/api/shared"
	"github.com/fernandodr19/library/pkg/instrumentation/logger"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Authorizer authorizes requests
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

// TrimSlashSuffix Removes the trailing slash from request, except if it is the root url.
func TrimSlashSuffix(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path != "/" {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	}

	next.ServeHTTP(w, r)
}

func AssureRequestID(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// create req id if none is provided
	reqID := w.Header().Get(shared.XReqID)
	if reqID == "" {
		reqID = uuid.NewString()
		w.Header().Set(shared.XReqID, reqID)
	}

	// insert log with req id on context
	log := logger.Default().WithField(shared.XReqID, reqID)
	ctx := logger.ToCtx(r.Context(), log)

	next.ServeHTTP(w, r.WithContext(ctx))
}

// Handle middleware function to treat rest responses.
func Handle(handler func(r *http.Request) responses.Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handles request
		response := handler(r)

		// post application precessing...
		log := logger.FromCtx(r.Context())

		if response.Error != nil {
			err := response.Error
			log.Error(err)
		}

		// Setting headers
		for key, value := range response.Headers() {
			w.Header().Set(key, value)
		}

		err := responses.SendJSON(w, response.Payload, response.Status)
		if err != nil {
			log.Error(err)
		}
	}
}
