package api

import (
	"net/http"

	library "github.com/fernandodr19/library/pkg"
	"github.com/fernandodr19/library/pkg/config"
	"github.com/fernandodr19/library/pkg/gateway/api/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	http_swagger "github.com/swaggo/http-swagger"
	"github.com/urfave/negroni"
)

// BuildHandler builds api handler
func BuildHandler(app *library.App, cfg *config.Config, auth middleware.Authorizer) (http.Handler, error) {
	r := mux.NewRouter()

	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods(http.MethodGet)
	r.PathPrefix("/healthcheck").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)
	r.PathPrefix("/docs/v1/library/swagger").Handler(http_swagger.WrapHandler).Methods(http.MethodGet)

	publicV1 := r.PathPrefix("/api/v1").Subrouter()
	adminV1 := r.PathPrefix("/admin/v1").Subrouter()
	accounts.NewHandler(publicV1, adminV1, app.Accounts, auth)

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n := negroni.New()
	n.UseFunc(middleware.TrimSlashSuffix)
	n.UseFunc(middleware.AssureRequestID)
	n.UseHandler(middleware.Cors(r))

	return n, nil
}
