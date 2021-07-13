package repositories

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	usecase "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/fernandodr19/library/pkg/gateway/repositories/sqlc"
	"github.com/jackc/pgx/v4"
)

var _ usecase.Repository = &AccountRepository{}

// AccountRepository is the repository of accounts
type AccountRepository struct {
	conn *pgx.Conn
	q    *sqlc.Queries
}

// NewAccountRepository builds an account repository
func NewAccountRepository(conn *pgx.Conn) *AccountRepository {
	return &AccountRepository{
		conn: conn,
		q:    sqlc.New(conn),
	}
}

// GetAccountByEmail gets an account for a given email
func (r AccountRepository) GetAccountByEmail(ctx context.Context, email vos.Email) (accounts.Account, error) {
	const operation = "repositories.AccountRepository.GetAccountByEmail"

	rawAcc, err := r.q.GetAccountByEmail(ctx, email.String())
	if err != nil {
		if err == pgx.ErrNoRows {
			return accounts.Account{}, usecase.ErrAccountNotFound
		}
		return accounts.Account{}, domain.Error(operation, err)
	}

	return mapRawAcc(rawAcc), nil
}

func mapRawAcc(a sqlc.Account) accounts.Account {
	return accounts.Account{
		ID:             vos.AccID(a.ID.String()),
		Email:          vos.Email(a.Email),
		HashedPassword: vos.HashedPassword(a.Password),
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

// CreateAccount creates an account on db
func (r AccountRepository) CreateAccount(ctx context.Context, email vos.Email, hashedPassword vos.HashedPassword) (vos.AccID, error) {
	const operation = "repositories.AccountRepository.CreateAccount"

	id, err := r.q.CreateAccount(ctx, sqlc.CreateAccountParams{
		Email:    email.String(),
		Password: hashedPassword.String(),
	})
	if err != nil {
		return "", domain.Error(operation, err)
	}

	return vos.AccID(id.String()), nil
}
