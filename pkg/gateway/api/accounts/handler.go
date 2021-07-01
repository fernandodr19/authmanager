package accounts

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/fernandodr19/library/pkg/gateway/auth"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usecase accounts.Usecase
}

func NewHandler(public *mux.Router, admin *mux.Router, usecase accounts.Usecase) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	a, err := auth.New("My Secret")
	if err != nil {
		panic(err)
	}

	fmt.Println(a.CreateToken("my-user", 30*time.Second))

	// public.Handle("/do-something", middleware.Handle(h.DoSomething)).Methods(http.MethodGet)
	public.Handle("/do-something",
		a.AuthorizeRequest(middleware.Handle(h.DoSomething))).
		Methods(http.MethodGet)

	return h
}

var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})
