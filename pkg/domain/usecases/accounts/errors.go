package accounts

import "errors"

var (
	ErrNotImplemented         = errors.New("not implemented")
	ErrAccountNotFound        = errors.New("account not found")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)
