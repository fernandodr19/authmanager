package repositories

import (
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
