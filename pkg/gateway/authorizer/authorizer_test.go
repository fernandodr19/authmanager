package authorizer

import (
	"testing"
	"time"

	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AuthToken(t *testing.T) {
	auth, err := New("my-secret-key")
	require.NoError(t, err)

	var userID vos.UserID = "user-id"
	t.Run("test authorizer happy path", func(t *testing.T) {
		accessToken, err := auth.CreateToken(userID, time.Minute)
		require.NoError(t, err)

		payload, err := auth.verifyToken(accessToken.String())
		require.NoError(t, err)

		assert.NotEmpty(t, payload.TokenID)
		assert.Equal(t, payload.UserID, userID)
		assert.Equal(t, payload.ExpiredAt.Sub(payload.IssuedAt), time.Minute)
	})

	t.Run("test authorizer key mismatch", func(t *testing.T) {
		accessToken, err := auth.CreateToken(userID, time.Minute)
		require.NoError(t, err)

		auth2, err := New("my-secret-different-key")
		require.NoError(t, err)

		_, err = auth2.verifyToken(accessToken.String())
		assert.Error(t, err) //ErrSignatureInvalid
	})

	t.Run("test authorizer expired token", func(t *testing.T) {
		accessToken, err := auth.CreateToken(userID, time.Nanosecond)
		require.NoError(t, err)

		time.Sleep(2 * time.Nanosecond)

		_, err = auth.verifyToken(accessToken.String())
		assert.Error(t, err) //ErrExpiredToken
	})
}
