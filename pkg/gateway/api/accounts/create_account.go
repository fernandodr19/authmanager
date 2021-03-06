package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/fernandodr19/authmanager/pkg/domain"
	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/fernandodr19/authmanager/pkg/gateway/api/responses"
)

// CreateAccount handles create account requests
// @Summary Creates an account
// @Description Creates an account for a given email.
// @Tags Accounts
// @Param Body body CreateAccountRequest true "Body"
// @Accept json
// @Produce json
// @Success 201 "Account successfully created"
// @Failure 400 "Could not parse request"
// @Failure 409 "User already registered"
// @Failure 422 "Request is well formed but contains invalid data"
// @Failure 500 "Internal server error"
// @Router /accounts [post]
func (h Handler) CreateAccount(r *http.Request) responses.Response {
	operation := "accounts.Handler.CreateAccount"

	ctx := r.Context()
	var body CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return responses.BadRequest(domain.Error(operation, err), responses.ErrInvalidBody)
	}

	_, err = h.Usecase.CreateAccount(ctx, body.Email, body.Password)
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.Created(nil)
}

// CreateAccountRequest payload
type CreateAccountRequest struct {
	Email    vos.Email    `json:"email"`
	Password vos.Password `json:"password"`
}
