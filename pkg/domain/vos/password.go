package vos

//nolint
type (
	Password       string
	HashedPassword string
)

// Valid checks if password is valid
func (p Password) Valid() bool {
	return p != ""
}
