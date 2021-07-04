package accounts

import "errors"

var (
	ErrNotImplemented         = errors.New("not implemented")
	ErrAccountNotFound        = errors.New("account not found")
	ErrWrongPassword          = errors.New("wrong password")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)
