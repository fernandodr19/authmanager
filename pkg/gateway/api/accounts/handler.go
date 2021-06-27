package accounts

import (
	"net/http"

	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usecase accounts.Usecase
}

func NewHandler(public *mux.Router, admin *mux.Router, usecase accounts.Usecase) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	public.HandleFunc("/do-something",
		middleware.Handle(h.DoSomething)).
		Methods(http.MethodGet)

	return h
}
