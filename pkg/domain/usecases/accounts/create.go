package accounts

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/instrumentation"
)

// CreateAccount creates a brand new account for a given user
func (u AccountsUsecase) CreateAccount(ctx context.Context, email vos.Email, password vos.Password) error {
	const operation = "accounts.AccountsUsecase.CreateAccount"

	// TODO: receiver encrypted params (maybe JWE)
	// instrumentation.Logger().WithField("TOKEN", ctx.Value(accounts.UserIDContextKey)).Info("sss")
	log := instrumentation.Logger().WithField("email", email)
	log.Infoln("creating account")

	if !email.Valid() {
		return domain.Error(operation, accounts.ErrInvalidEmail)
	}

	if !password.Valid() {
		return domain.Error(operation, accounts.ErrInvalidPassword)
	}

	// check if user is alreary registered
	_, err := u.AccountsRepository.GetAccountByEmail(ctx, email)
	if err != ErrAccountNotFound {
		return domain.Error(operation, ErrEmailAlreadyRegistered)
	}

	// hashes the password
	hashedPass, err := u.Encrypter.HashedPassword(password)
	if err != nil {
		return domain.Error(operation, err)
	}

	// create acc on db
	userID, err := u.AccountsRepository.CreateAccount(ctx, email, hashedPass)
	if err != nil {
		return domain.Error(operation, err)
	}

	log.WithField("useID", userID).Infoln("account created")

	return nil
}
