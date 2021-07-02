package vos

type (
	Email    string
	Password string
)

func (e Email) String() string {
	return string(e)
}

func (p Password) String() string {
	return string(p)
}
