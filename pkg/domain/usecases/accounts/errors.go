package accounts

import "errors"

var (
	ErrNotImplemented         = errors.New("not implemented")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)
