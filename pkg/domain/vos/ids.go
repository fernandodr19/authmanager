package vos

import "net/mail"

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

func (e Email) Valid() bool {
	_, err := mail.ParseAddress(e.String())
	return err == nil
}
