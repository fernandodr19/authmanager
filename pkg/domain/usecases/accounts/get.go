package accounts

import (
	"context"

	"github.com/fernandodr19/authmanager/pkg/domain"
	"github.com/fernandodr19/authmanager/pkg/domain/entities/accounts"
	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/fernandodr19/authmanager/pkg/instrumentation/logger"
)

// GetAccountDetaiils retrieves an account for a given user
func (u AccountsUsecase) GetAccountDetaiils(ctx context.Context, accID vos.AccID) (accounts.Account, error) {
	const operation = "accounts.AccountsUsecase.GetAccountDetaiils"

	userID := ctx.Value(accounts.UserIDContextKey)
	log := logger.FromCtx(ctx).WithField("accID", userID)

	log.Info("get account details")

	acc, err := u.AccountsRepository.GetAccountByID(ctx, accID)
	if err != nil {
		return accounts.Account{}, domain.Error(operation, err)
	}

	log.Info("got account details")

	return acc, nil
}
