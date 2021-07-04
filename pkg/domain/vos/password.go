package vos

type (
	Password       string
	HashedPassword string
)

func (p Password) Valid() bool {
	return p != ""
}
