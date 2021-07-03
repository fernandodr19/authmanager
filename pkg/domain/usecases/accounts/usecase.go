package accounts

import (
	"context"
	"time"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/instrumentation"
	"github.com/sirupsen/logrus"
)

//go:generate moq -skip-ensure -stub -out mocks.gen.go . Usecase:AccountsMockUsecase

type TokenGenerator interface {
	CreateTokens(userID vos.UserID, accessDuration time.Duration, refreshDuration time.Duration) (vos.Tokens, error)
}

var _ Usecase = &AccountsUsecase{}

type Usecase interface {
	CreateAccount(context.Context, vos.Email, vos.Password) error
}

type Repository interface {
	GetAccountByEmail(context.Context, vos.Email) (accounts.Account, error)
	CreateAccount(context.Context, vos.Email, vos.HashedPassword) (vos.UserID, error)
	Login(context.Context) error
	Logout(context.Context) error
}

type AccountsUsecase struct {
	AccountsRepository Repository
	TokenGenerator     TokenGenerator
}

func NewAccountsUsecase(accRepo Repository, tokenGenerator TokenGenerator) *AccountsUsecase {
	return &AccountsUsecase{
		AccountsRepository: accRepo,
	}
}

func (a AccountsUsecase) CreateAccount(ctx context.Context, email vos.Email, password vos.Password) error {
	const operation = "accounts.AccountsUsecase.CreateAccount"
	// instrumentation.Logger().WithField("TOKEN", ctx.Value(accounts.UserIDContextKey)).Info("sss")
	instrumentation.Logger().WithField("email", email).Infoln("Creating account")

	if !email.Valid() {
		return domain.Error(operation, accounts.ErrInvalidEmail)
	}

	if !password.Valid() {
		return domain.Error(operation, accounts.ErrInvalidPassword)
	}

	// check if user is alreary registered
	_, err := a.AccountsRepository.GetAccountByEmail(ctx, email)
	if err != ErrAccountNotFound {
		return domain.Error(operation, ErrEmailAlreadyRegistered)
	}

	// hashes the password
	hashedPass, err := password.Hashed()
	if err != nil {
		return domain.Error(operation, err)
	}

	// create acc on db
	userID, err := a.AccountsRepository.CreateAccount(ctx, email, hashedPass)
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
