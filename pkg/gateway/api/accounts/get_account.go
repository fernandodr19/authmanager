package accounts

import (
	"net/http"
	"time"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetAccount handles get account requests
// @Summary Get account details
// @Description Get account details for a given ID.
// @Tags Accounts
// @Param Authorization header string true "Bearer Authorization Token"
// @Param account_id path string true "Account ID"
// @Accept json
// @Produce json
// @Success 200 {object} GetAccountResponse
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /accounts/{account_id} [get]
func (h Handler) GetAccount(r *http.Request) responses.Response {
	operation := "accounts.Handler.GetAccount"

	ctx := r.Context()

	accID, err := uuid.Parse(mux.Vars(r)["account_id"])
	if err != nil {
		return responses.BadRequest(domain.Error(operation, err), responses.ErrInvalidUserID)
	}

	acc, err := h.Usecase.GetAccountDetaiils(ctx, vos.AccID(accID.String()))
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	return responses.OK(GetAccountResponse{
		AccountID: acc.ID,
		Email:     acc.Email,
		CratedAt:  acc.CreatedAt,
	})
}

type GetAccountResponse struct {
	AccountID vos.AccID `json:"account_id"`
	Email     vos.Email `json:"email"`
	CratedAt  time.Time `json:"created_at"`
}
