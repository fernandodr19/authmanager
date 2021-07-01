package accounts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/fernandodr19/library/pkg/gateway/authorizer"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usecase accounts.Usecase
}

func NewHandler(public *mux.Router, admin *mux.Router, usecase accounts.Usecase) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	auth, err := authorizer.New("My Secret")
	if err != nil {
		panic(err)
	}

	fmt.Println(auth.CreateToken("my-user", 3000*time.Second))

	// public.Handle("/do-something", middleware.Handle(h.DoSomething)).Methods(http.MethodGet)
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
