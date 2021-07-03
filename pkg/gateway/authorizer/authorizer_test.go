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
		tokens, err := auth.CreateTokens(userID, time.Minute, time.Hour)
		require.NoError(t, err)

		accessPayload, err := auth.verifyToken(tokens.AccessToken.String())
		require.NoError(t, err)

		assert.NotEmpty(t, accessPayload.TokenID)
		assert.Equal(t, accessPayload.UserID, userID)
		assert.Equal(t, accessPayload.ExpiredAt.Sub(accessPayload.IssuedAt), time.Minute)

		refreshPayload, err := auth.verifyToken(tokens.RefreshToken.String())
		require.NoError(t, err)

		assert.NotEmpty(t, refreshPayload.TokenID)
		assert.Equal(t, refreshPayload.UserID, userID)
		assert.Equal(t, refreshPayload.ExpiredAt.Sub(refreshPayload.IssuedAt), time.Hour)
	})

	t.Run("test authorizer key mismatch", func(t *testing.T) {
		tokens, err := auth.CreateTokens(userID, time.Minute, time.Hour)
		require.NoError(t, err)

		auth2, err := New("my-secret-different-key")
		require.NoError(t, err)

		_, err = auth2.verifyToken(tokens.AccessToken.String())
		assert.Error(t, err) //ErrSignatureInvalid
	})

	t.Run("test authorizer expired token", func(t *testing.T) {
		accessToken, err := auth.CreateTokens(userID, time.Nanosecond, time.Hour)
		require.NoError(t, err)

		time.Sleep(2 * time.Nanosecond)

		_, err = auth.verifyToken(accessToken.AccessToken.String())
		assert.Error(t, err) //ErrExpiredToken
	})
}
