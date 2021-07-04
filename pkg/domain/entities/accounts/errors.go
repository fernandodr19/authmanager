package accounts

import "errors"

var (
	ErrInvalidEmail    = errors.New("invalid email")    // invalid email
	ErrInvalidPassword = errors.New("invalid password") // invalid password
)
