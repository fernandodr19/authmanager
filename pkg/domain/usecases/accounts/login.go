package accounts

import (
	"context"
	"time"

	"github.com/fernandodr19/authmanager/pkg/domain"
	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/fernandodr19/authmanager/pkg/instrumentation/logger"
)

// Login authenticates the user
func (u AccountsUsecase) Login(ctx context.Context, email vos.Email, password vos.Password) (vos.AccID, vos.Tokens, error) {
	const operation = "accounts.AccountsUsecase.Login"

	// TODO: receiver encrypted params (maybe JWE)
	log := logger.FromCtx(ctx)
	log.WithField("email", email).Infoln("user login attempt")

	tokens := vos.Tokens{}

	// retrieve user from db
	acc, err := u.AccountsRepository.GetAccountByEmail(ctx, email)
	if err != nil {
		if err == ErrAccountNotFound {
			return "", tokens, domain.Error(operation, ErrAccountNotFound)
		}
		return "", tokens, domain.Error(operation, err)
	}

	// check password
	if !u.Encrypter.PasswordMathces(password, acc.HashedPassword) {
		return "", tokens, domain.Error(operation, ErrWrongPassword)
	}

	// generate tokens
	tokens, err = u.TokenGenerator.CreateTokens(acc, 10*time.Minute, 2*time.Hour)
	if err != nil {
		return "", tokens, domain.Error(operation, err)
	}

	log.WithField("email", email).Infoln("user logged in successfully")

	return acc.ID, tokens, nil
}
