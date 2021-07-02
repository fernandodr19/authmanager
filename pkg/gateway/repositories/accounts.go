package repositories

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	acc_usecase "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AccountRepository struct {
	*pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		Pool: db,
	}
}

func (a AccountRepository) GetAccountByEmail(ctx context.Context, email vos.Email) (accounts.Account, error) {
	return accounts.Account{
		Email: email,
	}, nil
}

func (a AccountRepository) CreateAccount(context.Context) error {
	return acc_usecase.ErrNotImplemented
}

func (a AccountRepository) Login(context.Context) error {
	return acc_usecase.ErrNotImplemented
}

func (a AccountRepository) Logout(context.Context) error {
	return acc_usecase.ErrNotImplemented
}
