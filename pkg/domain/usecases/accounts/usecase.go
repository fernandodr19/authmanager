package accounts

import (
	"context"
	"net/mail"
	"time"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/instrumentation"
)

//go:generate moq -skip-ensure -stub -out mocks.gen.go . Usecase:AccountsMockUsecase

type Usecase interface {
	CreateAccount(context.Context, vos.Email, vos.Password) (Tokens, error)
}

type TokenGenerator interface {
	CreateToken(userID vos.UserID, duration time.Duration) (vos.AccessToken, error)
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

func (a AccountsUsecase) CreateAccount(ctx context.Context, email vos.Email, password vos.Password) (Tokens, error) {
	const operation = "accounts.AccountsUsecase.CreateAccount"
	// instrumentation.Logger().WithField("TOKEN", ctx.Value(accounts.UserIDContextKey)).Info("sss")
	instrumentation.Logger().WithField("email", email).Infoln("Creating account")

	if !validateEmail(email) {
		return Tokens{}, domain.Error(operation, accounts.ErrInvalidEmail)
	}

	if !validatePassword(password) {
		return Tokens{}, domain.Error(operation, accounts.ErrInvalidPassword)
	}

	_, err := a.AccountsRepository.GetAccountByEmail(ctx, email)
	if err != nil {
		// SHOULD RESULT IN STATUS CONFLICT, TEST IT
		return Tokens{}, domain.Error(operation, ErrAlreadyRegistered)
	}

	// creates on db
	a.AccountsRepository.CreateAccount(ctx)

	// generates token ONLY on login, change that

	return Tokens{}, nil
}

type Tokens struct {
	AccessToken  vos.AccessToken
	RefreshToken vos.RefreshToken
}

func validateEmail(email vos.Email) bool {
	_, err := mail.ParseAddress(email.String())
	return err != nil
}

func validatePassword(password vos.Password) bool {
	return password != ""
}
