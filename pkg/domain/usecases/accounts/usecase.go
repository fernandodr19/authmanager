package accounts

import (
	"context"

	"github.com/fernandodr19/authenticator/pkg/domain/entities/accounts"
)

// Maybe move interface for where it is gonna be used
type Usecase interface {
	DoSomething(context.Context) error
}

var _ Usecase = &AccountsUsecase{}

type AccountsUsecase struct {
	AccountsRepository accounts.Repository
}

func NewAccountsUsecase(accRepo accounts.Repository) *AccountsUsecase {
	return &AccountsUsecase{
		AccountsRepository: accRepo,
	}
}

func (a AccountsUsecase) DoSomething(ctx context.Context) error {
	return ErrNotImplemented
}
