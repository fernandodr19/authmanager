package encrypter

import (
	"github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"golang.org/x/crypto/bcrypt"
)

var _ accounts.Encrypter = &Encrypter{}

type Encrypter struct {
}

func (c Encrypter) HashedPassword(password vos.Password) (vos.HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return vos.HashedPassword(hash), nil
}

func (c Encrypter) PasswordMathces(password vos.Password, hashedPassword vos.HashedPassword) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
