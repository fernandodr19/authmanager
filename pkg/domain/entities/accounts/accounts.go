package accounts

import (
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/google/uuid"
)

const UserIDContextKey vos.UserID = "user-id-context-key"

type Account struct {
	ID    uuid.UUID
	Email vos.Email
}
