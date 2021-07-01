package repositories

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
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

func (a AccountRepository) SignUp(context.Context) error {
	return accounts.ErrNotImplemented
}

func (a AccountRepository) Login(context.Context) error {
	return accounts.ErrNotImplemented
}

func (a AccountRepository) Logout(context.Context) error {
	return accounts.ErrNotImplemented
}
