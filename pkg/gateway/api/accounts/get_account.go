package accounts

import (
	"fmt"
	"net/http"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/gateway/api/responses"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Login
// @Summary Authenticate user credentials
// @Description Authenticate user credentials.
// @Tags Accounts
// @Param Authorization header string true "Bearer Authorization Token"
// @Param acc_id path string true "Account ID"
// @Accept json
// @Produce json
// @Success 200 {object} LoginResponse
// @Failure 401 "Invalid password"
// @Failure 404 "User not found"
// @Failure 500 "Internal server error"
// @Router /accounts/{acc_id} [get]
// Login handles login requests
func (h Handler) GetAccount(r *http.Request) responses.Response {
	operation := "accounts.Handler.GetAccount"

	ctx := r.Context()

	accID, err := uuid.Parse(mux.Vars(r)["acc_id"])
	if err != nil {
		return responses.BadRequest(domain.Error(operation, err), responses.ErrInvalidUserID)
	}

	fmt.Println(accID)

	acc, err := h.Usecase.GetAccountDetaiils(ctx)
	if err != nil {
		return responses.ErrorResponse(domain.Error(operation, err))
	}

	fmt.Print(acc)

	return responses.OK(nil)
}
