package api

import (
	"net/http"

	authenticator "github.com/fernandodr19/authenticator/pkg"
	"github.com/fernandodr19/authenticator/pkg/config"
	"github.com/fernandodr19/authenticator/pkg/gateway/api/accounts"
	"github.com/fernandodr19/authenticator/pkg/gateway/api/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	http_swagger "github.com/swaggo/http-swagger"
	"github.com/urfave/negroni"
)

func BuildHandler(app *authenticator.App, cfg *config.Config) (http.Handler, error) {
	r := mux.NewRouter()

	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods(http.MethodGet)
	r.PathPrefix("/healthcheck").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)
	r.PathPrefix("/docs/v1/authenticator/swagger").Handler(http_swagger.WrapHandler).Methods(http.MethodGet)

	publicV1 := r.PathPrefix("/api/v1").Subrouter()
	adminV1 := r.PathPrefix("/admin/v1").Subrouter()
	accounts.NewHandler(publicV1, adminV1, app.Accounts)

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n := negroni.New()
	n.UseFunc(middleware.TrimSlashSuffix)
	n.UseHandler(middleware.Cors(r))

	return n, nil
}
