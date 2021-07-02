package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/fernandodr19/library/pkg/instrumentation"
)

// Create account
// @Summary Creates an account
// @Description Creates an account for a given email.
// @Tags Accounts
// @Param Body body CreateAccountRequest true "Body"
// @Param Authorization header string true "Bearer Authorization Token"
// @Accept json
// @Produce json
// @Success 201 {object} CreateAccountResponse
// @Header 201 {string} Token "X-Request-Id"
// @Failure 500 "Internal server error"
// @Router /signup [post]
func (h Handler) CreateAccount(r *http.Request) responses.Response {
	operation := "accounts.Handler.CreateAccount"

	ctx := r.Context()
	log := instrumentation.Logger()
	log.Info("doing somethings")

	var body CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return responses.BadRequest(domain.Error(operation, err), responses.ErrInvalidBody)
	}

	tokesn, err := h.Usecase.CreateAccount(ctx, body.Email, body.Password)
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.OK(CreateAccountResponse{
		AccessToken:  tokesn.AccessToken,
		RefreshToken: tokesn.RefreshToken,
	})
}

type CreateAccountRequest struct {
	Email    vos.Email    `json:"email"`
	Password vos.Password `json:"password"`
}

type CreateAccountResponse struct {
	AccessToken  vos.AccessToken  `json:"access_token"`
	RefreshToken vos.RefreshToken `json:"refresh_token"`
}
