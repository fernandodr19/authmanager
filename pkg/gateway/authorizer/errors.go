package authorizer

import "errors"

var (
	// ErrInvalidToken invalid token error
	ErrInvalidToken = errors.New("token is invalid")
	// ErrExpiredToken expired token error
	ErrExpiredToken = errors.New("token has expired")
)
