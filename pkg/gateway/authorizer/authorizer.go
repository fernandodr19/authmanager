package authorizer

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fernandodr19/library/pkg/domain"
	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	acc_usecase "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/library/pkg/domain/vos"
	"github.com/google/uuid"
)

var _ acc_usecase.TokenGenerator = &bearerAuthorizer{}

type bearerAuthorizer struct {
	secretKey []byte
}

// New is the bearer auhtorizer builder
func New(secretKey string) (*bearerAuthorizer, error) {
	return &bearerAuthorizer{[]byte(secretKey)}, nil
}

type payload struct {
	TokenID   uuid.UUID  `json:"token_id"`
	UserID    vos.UserID `json:"user_id"`
	IssuedAt  time.Time  `json:"issued_at"`
	ExpiredAt time.Time  `json:"expired_at"`
}

// Valid checks if a payload is valid
func (p *payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// CreateTokens generate both access & refresh tokens
func (b *bearerAuthorizer) CreateTokens(userID vos.UserID, accessDuration time.Duration, refreshDuration time.Duration) (vos.Tokens, error) {
	const operation = "authorizer.BearerAuthorizer.CreateToken"

	accessToken, err := b.createToken(userID, accessDuration)
	if err != nil {
		return vos.Tokens{}, domain.Error(operation, err)
	}

	refreshToken, err := b.createToken(userID, refreshDuration)
	if err != nil {
		return vos.Tokens{}, domain.Error(operation, err)
	}

	return vos.Tokens{
		AccessToken:  vos.AccessToken(accessToken),
		RefreshToken: vos.RefreshToken(refreshToken),
	}, nil
}

func (b *bearerAuthorizer) createToken(userID vos.UserID, duration time.Duration) (string, error) {
	const operation = "authorizer.createAccessToken"

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", domain.Error(operation, err)
	}

	now := time.Now()

	p := &payload{
		TokenID:   tokenID,
		UserID:    userID,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, p)
	token, err := jwtToken.SignedString(b.secretKey)
	if err != nil {
		return "", domain.Error(operation, err)
	}

	return token, nil
}

// AuthorizeRequest is a middleware that handles request authorization
func (b *bearerAuthorizer) AuthorizeRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// TODO use responses package maybe?
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing auth header"))
			return
		}

		splitedAuthHeader := strings.Split(authHeader, " ")
		if len(splitedAuthHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid auth header"))
			return
		}

		token := splitedAuthHeader[1]

		payload, err := b.verifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), accounts.UserIDContextKey, payload.UserID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (b *bearerAuthorizer) verifyToken(token string) (*payload, error) {
	const operation = "authorizer.BearerAuthorizer.VerifyToken"

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, domain.Error(operation, ErrInvalidToken)
		}
		return b.secretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &payload{}, keyFunc)
	if err != nil {
		return nil, domain.Error(operation, err)
	}

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, domain.Error(operation, ErrExpiredToken)
		}
		return nil, domain.Error(operation, ErrInvalidToken)
	}

	p, ok := jwtToken.Claims.(*payload)
	if !ok {
		return nil, domain.Error(operation, ErrInvalidToken)
	}

	return p, nil
}
