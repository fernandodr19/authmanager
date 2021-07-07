package vos

//nolint
type (
	Password       string
	HashedPassword string
)

// String returns password as string
func (p Password) String() string {
	return string(p)
}

// Valid checks if password is valid
func (p Password) Valid() bool {
	return p != ""
}

// String returns hashed password as string
func (h HashedPassword) String() string {
	return string(h)
}
