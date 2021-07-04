package authorizer

import "errors"

var (
	// ErrInvalidToken invalid token error
	ErrInvalidToken = errors.New("token is invalid")
	// ErrInvalidToken expired token error
	ErrExpiredToken = errors.New("token has expired")
)
