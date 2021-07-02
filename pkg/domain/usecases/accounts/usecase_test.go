package accounts

import (
	"testing"

	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/stretchr/testify/assert"
)

func Test_ValidEmail(t *testing.T) {
	testTable := []struct {
		Name  string
		Email vos.Email
		Error bool
	}{
		{
			Name:  "valid email",
			Email: "valid@gmail.com",
			Error: false,
		},
		{
			Name:  "valid email with plus sign",
			Email: "valid+123@gmail.com",
			Error: false,
		},
		{
			Name:  "invalid without @",
			Email: "invalid",
			Error: true,
		},
		{
			Name:  "invalid without domain",
			Email: "invalid@",
			Error: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Error, validateEmail(tt.Email))
		})
	}
}
