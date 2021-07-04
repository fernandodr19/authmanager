package accounts

import (
	"context"
	"time"

	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
)

//go:generate moq -skip-ensure -stub -out mocks.gen.go . Usecase:AccountsMockUsecase

var _ Usecase = &AccountsUsecase{}

// Accounts usecase
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
	Encrypter          Encrypter
}

func NewAccountsUsecase(accRepo Repository, tokenGenerator TokenGenerator, encrypter Encrypter) *AccountsUsecase {
	return &AccountsUsecase{
		AccountsRepository: accRepo,
		Encrypter:          encrypter,
	}
}

type TokenGenerator interface {
	CreateTokens(userID vos.UserID, accessDuration time.Duration, refreshDuration time.Duration) (vos.Tokens, error)
}

type Encrypter interface {
	HashedPassword(password vos.Password) (vos.HashedPassword, error)
	PasswordMathces(password vos.Password, hashedPassword vos.HashedPassword) bool
}
