package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Cors applies cors rules to router
func Cors(r *mux.Router) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Origin", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	return handlers.CORS(originsOk, headersOk, methodsOk)(r)
}

// Removes the trailing slash from request, except if it is the root url.
// If the url is https://stone.com.br/api or https://stone.com.br/api/
// both will match.
// This was done as gorilla mux default method for this doesn't support POST requests: https://github.com/gorilla/mux/issues/30
// Usage:
// n := negroni.Classic()
// n.UseFunc(middleware.TrimSlashSuffix)
func TrimSlashSuffix(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Path != "/" {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
	}

	next.ServeHTTP(w, r)
}
