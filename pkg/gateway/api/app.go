package api

import (
	"net/http"

	"github.com/fernandodr19/authenticator/pkg/config"
	"github.com/fernandodr19/authenticator/pkg/gateway/api/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

func BuildHandler(cfg *config.Config) (http.Handler, error) {
	r := mux.NewRouter()

	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods(http.MethodGet)
	r.PathPrefix("/healthcheck").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	n := negroni.New()
	n.UseFunc(middleware.TrimSlashSuffix)
	n.UseHandler(middleware.Cors(r))

	return n, nil
}
