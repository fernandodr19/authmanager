package accounts

import "errors"

var (
	ErrNotImplemented    = errors.New("not implemented")
	ErrAlreadyRegistered = errors.New("email already registered")
)
