package encrypter

import (
	"testing"

	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PasswordEncryption(t *testing.T) {
	c := Encrypter{}

	t.Run("password matches", func(t *testing.T) {
		p := vos.Password("mypass")
		hash, err := c.HashedPassword(p)
		require.NoError(t, err)
		assert.True(t, c.PasswordMathces(p, hash))
	})

	t.Run("password does not match", func(t *testing.T) {
		p := vos.Password("mypass")
		hash, err := c.HashedPassword(p)
		require.NoError(t, err)
		p2 := vos.Password("wrongpass")
		assert.False(t, c.PasswordMathces(p2, hash))
	})

}
