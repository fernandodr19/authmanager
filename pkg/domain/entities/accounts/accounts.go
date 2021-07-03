package accounts

import (
	"time"

	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/google/uuid"
)

const UserIDContextKey vos.UserID = "user-id-context-key"

type Account struct {
	ID        uuid.UUID
	Email     vos.Email
	Password  vos.Password
	CreatedAt time.Time
	UpdatedAt time.Time
}
