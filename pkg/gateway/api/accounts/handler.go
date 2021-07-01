package accounts

import (
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/gateway/api/middleware"
	"github.com/form3tech-oss/jwt-go"

	"github.com/gorilla/mux"
)

type Handler struct {
	Usecase accounts.Usecase
}

func NewHandler(public *mux.Router, admin *mux.Router, usecase accounts.Usecase) *Handler {
	h := &Handler{
		Usecase: usecase,
	}

	// token := jwt.Token{
	// 	Raw:       "",
	// 	Method:    nil,
	// 	Header:    map[string]interface{}{},
	// 	Claims:    nil,
	// 	Signature: "",
	// 	Valid:     false,
	// }

	auth := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	// public.Handle("/do-something", middleware.Handle(h.DoSomething)).Methods(http.MethodGet)
	public.Handle("/do-something",
		auth.Handler(middleware.Handle(h.DoSomething))).
		Methods(http.MethodGet)

	return h
}

var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

})
