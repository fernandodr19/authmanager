package repositories

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	usecase "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/jackc/pgx/v4"
)

var _ usecase.Repository = &AccountRepository{}

// AccountRepository is the repository of accounts
type AccountRepository struct {
	Conn *pgx.Conn
}

// NewAccountRepository builds an account repository
func NewAccountRepository(db *pgx.Conn) *AccountRepository {
	return &AccountRepository{
		Conn: db,
	}
}

// GetAccountByEmail gets an account for a given email
func (r AccountRepository) GetAccountByEmail(ctx context.Context, email vos.Email) (accounts.Account, error) {
	const operation = "repositories.AccountRepository.GetAccountByEmail"

	const cmd = `
		SELECT
			id,
			email,
			password,
			created_at,
			updated_at
		FROM accounts
		WHERE email = $1
	`
	var acc accounts.Account
	err := r.Conn.QueryRow(ctx, cmd, email).
		Scan(&acc.ID,
			&acc.Email,
			&acc.HashedPassword,
			&acc.CreatedAt,
			&acc.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return acc, usecase.ErrAccountNotFound
		}
		return acc, domain.Error(operation, err)
	}

	return acc, nil
}

// CreateAccount creates an account on db
func (r AccountRepository) CreateAccount(ctx context.Context, email vos.Email, hashedPassword vos.HashedPassword) (vos.UserID, error) {
	const operation = "repositories.AccountRepository.CreateAccount"

	const cmd = `
		INSERT INTO accounts (email, password)
		VALUES ($1, $2)
		RETURNING id
	`
	var userID vos.UserID
	err := r.Conn.QueryRow(ctx, cmd, email, hashedPassword).
		Scan(&userID)
	if err != nil {
		return "", domain.Error(operation, err)
	}

	return userID, nil
}
