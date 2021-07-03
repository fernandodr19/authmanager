package repositories

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	acc_usecase "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type AccountRepository struct {
	Conn *pgx.Conn
}

func NewAccountRepository(db *pgx.Conn) *AccountRepository {
	return &AccountRepository{
		Conn: db,
	}
}

func (a AccountRepository) GetAccountByEmail(ctx context.Context, email vos.Email) (accounts.Account, error) {
	return accounts.Account{
		Email: email,
	}, nil
}

func (a AccountRepository) CreateAccount(context.Context) (vos.UserID, error) {
	return vos.UserID(uuid.NewString()), nil
}

func (a AccountRepository) Login(context.Context) error {
	return acc_usecase.ErrNotImplemented
}

func (a AccountRepository) Logout(context.Context) error {
	return acc_usecase.ErrNotImplemented
}
