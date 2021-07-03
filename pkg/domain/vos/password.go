package vos

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	Password       string
	HashedPassword string
)

func (p Password) Valid() bool {
	return p != ""
}

// TODO: move encrypt implementation out of domain
func (p Password) Hashed() (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return HashedPassword(hash), nil
}

func (h HashedPassword) Mathces(password Password) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(password))
	return err == nil
}
