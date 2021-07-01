package accounts

import "context"

type Repository interface {
	SignUp(context.Context) error
	Login(context.Context) error
	Logout(context.Context) error
}
