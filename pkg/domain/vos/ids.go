package vos

import (
	"net/mail"
)

//nolint
type (
	AccID string
	Email string
)

func (id AccID) String() string {
	return string(id)
}

func (e Email) String() string {
	return string(e)
}

// Valid validates if the email is valid
func (e Email) Valid() bool {
	_, err := mail.ParseAddress(e.String())
	return err == nil
}
