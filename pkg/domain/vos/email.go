package vos

type (
	Email string
)

func (e Email) String() string {
	return string(e)
}
