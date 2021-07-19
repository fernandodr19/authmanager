package accounts

import (
	"net/http"

	"github.com/fernandodr19/authmanager/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/authmanager/pkg/gateway/api/middleware"

	"github.com/gorilla/mux"
)

// Handler handles account related requests
type Handler struct {
	Usecase accounts.Usecase
}

// NewHandler builds accounts handler
func NewHandler(public *mux.Router, admin *mux.Router, usecase accounts.Usecase, auth middleware.Authorizer) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	public.Handle("/accounts",
		middleware.Handle(h.CreateAccount)).
		Methods(http.MethodPost)

	public.Handle("/accounts/login",
		middleware.Handle(h.Login)).
		Methods(http.MethodPost)

	public.Handle("/accounts/{account_id}",
		auth.AuthorizeRequest(middleware.Handle(h.GetAccount))).
		Methods(http.MethodGet)

	return h
}
