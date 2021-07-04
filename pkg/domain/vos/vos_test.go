package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EmailValidation(t *testing.T) {
	testTable := []struct {
		Name  string
		Email Email
		Valid bool
	}{
		{
			Name:  "valid email",
			Email: "valid@gmail.com",
			Valid: true,
		},
		{
			Name:  "valid email with plus sign",
			Email: "valid+123@gmail.com",
			Valid: true,
		},
		{
			Name:  "invalid without @",
			Email: "invalid",
			Valid: false,
		},
		{
			Name:  "invalid without domain",
			Email: "invalid@",
			Valid: false,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Valid, tt.Email.Valid())
		})
	}
}
