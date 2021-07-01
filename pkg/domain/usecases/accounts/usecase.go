package accounts

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
)

//go:generate moq -skip-ensure -stub -out mocks.gen.go . Usecase:AccountsMockUsecase

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
	a.AccountsRepository.SignUp(ctx)
	return nil
}
