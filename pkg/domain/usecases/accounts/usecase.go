package accounts

import (
	"context"
	"time"

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
	// instrumentation.Logger().WithField("TOKEN", ctx.Value(accounts.UserIDContextKey)).Info("sss")
	instrumentation.Logger().WithField("email", email).Infoln("Creating account")

	// validate email valid

	// validates email already registered

	// validate password is secure

	// creates on db
	a.AccountsRepository.SignUp(ctx)

	// generates token

	return Tokens{}, nil
}

type Tokens struct {
	AccessToken  vos.AccessToken
	RefreshToken vos.RefreshToken
}
