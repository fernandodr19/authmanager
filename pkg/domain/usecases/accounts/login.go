package accounts

import (
	"context"
	"time"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/instrumentation"
)

// Login authenticates the user
func (u AccountsUsecase) Login(ctx context.Context, email vos.Email, password vos.Password) (vos.Tokens, error) {
	const operation = "accounts.AccountsUsecase.Login"

	// TODO: receiver encrypted params (maybe JWE)
	instrumentation.Logger().WithField("email", email).Infoln("User login attempt")

	tokens := vos.Tokens{}

	// retrieve user from db
	acc, err := u.AccountsRepository.GetAccountByEmail(ctx, email)
	if err != nil {
		return tokens, domain.Error(operation, ErrAccountNotFound)
	}

	// check password
	if !u.Encrypter.PasswordMathces(password, acc.HashedPassword) {
		return tokens, domain.Error(operation, ErrInvalidPassword)
	}

	tokens, err = u.TokenGenerator.CreateTokens(acc, 10*time.Minute, 2*time.Hour)
	if err != nil {
		return tokens, domain.Error(operation, err)
	}

	return tokens, ErrNotImplemented
}
