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

func NewHandler(public *mux.Router, admin *mux.Router, usecase accounts.Usecase, auth middleware.Authorizer) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	public.Handle("/do-something",
		middleware.Handle(h.DoSomething)).
		Methods(http.MethodGet)

	public.Handle("/do-something-auth",
		auth.AuthorizeRequest(middleware.Handle(h.DoSomething))).
		Methods(http.MethodGet)

	return h
}

var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})
