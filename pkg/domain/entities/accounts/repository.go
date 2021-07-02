package accounts

import (
	"context"

	"github.com/fernandodr19/library/pkg/domain/vos"
)

type Repository interface {
	GetAccountByEmail(context.Context, vos.Email) (Account, error)
	CreateAccount(context.Context) error
	Login(context.Context) error
	Logout(context.Context) error
}
