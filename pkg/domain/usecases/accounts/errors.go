package accounts

import "errors"

var (
	ErrNotImplemented         = errors.New("not implemented")
	ErrAccountNotFound        = errors.New("account not found")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)
