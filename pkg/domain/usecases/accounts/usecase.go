package accounts

import (
	"context"
	"net/mail"
	"time"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/instrumentation"
	"github.com/sirupsen/logrus"
)

//go:generate moq -skip-ensure -stub -out mocks.gen.go . Usecase:AccountsMockUsecase

type Usecase interface {
	CreateAccount(context.Context, vos.Email, vos.Password) error
}

type TokenGenerator interface {
	CreateTokens(userID vos.UserID, accessDuration time.Duration, refreshDuration time.Duration) (vos.Tokens, error)
}

var _ Usecase = &AccountsUsecase{}

type AccountsUsecase struct {
	AccountsRepository accounts.Repository
	TokenGenerator     TokenGenerator
}

func NewAccountsUsecase(accRepo accounts.Repository, tokenGenerator TokenGenerator) *AccountsUsecase {
	return &AccountsUsecase{
		AccountsRepository: accRepo,
	}
}

func (a AccountsUsecase) CreateAccount(ctx context.Context, email vos.Email, password vos.Password) error {
	const operation = "accounts.AccountsUsecase.CreateAccount"
	// instrumentation.Logger().WithField("TOKEN", ctx.Value(accounts.UserIDContextKey)).Info("sss")
	instrumentation.Logger().WithField("email", email).Infoln("Creating account")

	if !validEmail(email) {
		return domain.Error(operation, accounts.ErrInvalidEmail)
	}

	if !validPassword(password) {
		return domain.Error(operation, accounts.ErrInvalidPassword)
	}

	_, err := a.AccountsRepository.GetAccountByEmail(ctx, email)
	if err != nil {
		// SHOULD RESULT IN STATUS CONFLICT, TEST IT
		return domain.Error(operation, ErrEmailAlreadyRegistered)
	}

	// creates on db hashed pass
	userID, err := a.AccountsRepository.CreateAccount(ctx)
	if err != nil {
		return domain.Error(operation, err)
	}

	// generates token ONLY on login, change that
	// tokens, err := a.TokenGenerator.CreateTokens(userID, 5*time.Minute, 2*time.Hour)
	// if err != nil {
	// 	return domain.Error(operation, err)
	// }

	instrumentation.Logger().WithFields(logrus.Fields{
		"email":  email,
		"userID": userID,
	}).Infoln("Account created")

	return nil
}

func validEmail(email vos.Email) bool {
	_, err := mail.ParseAddress(email.String())
	return err == nil
}

func validPassword(password vos.Password) bool {
	return password != ""
}
