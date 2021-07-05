package accounts

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/instrumentation/logger"
)

// GetAccountDetaiils retrieves an account for a given user
func (u AccountsUsecase) GetAccountDetaiils(ctx context.Context) (accounts.Account, error) {
	const operation = "accounts.AccountsUsecase.GetAccountDetaiils"

	userID := ctx.Value(accounts.UserIDContextKey)
	log := logger.FromCtx(ctx).WithField("userID", userID)

	log.Info("get account details")

	return accounts.Account{}, domain.Error(operation, ErrNotImplemented)
}
