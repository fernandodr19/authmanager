package vos

import "net/mail"

//nolint
type (
	UserID string
	Email  string
)

func (id UserID) String() string {
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
