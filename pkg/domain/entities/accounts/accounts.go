package accounts

import (
	"time"

	"github.com/fernandodr19/library/pkg/domain/vos"
)

const UserIDContextKey vos.AccID = "user-id-context-key"

type Account struct {
	ID             vos.AccID
	Email          vos.Email
	Password       vos.Password
	HashedPassword vos.HashedPassword
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
