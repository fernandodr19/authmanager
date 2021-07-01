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
	"github.com/fernandodr19/library/pkg/instrumentation"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type ContextKey string

const ContextKeyToken ContextKey = "token"

type Authorizer interface {
	CreateToken(username string, duration time.Duration) (string, error)
	AuthorizeRequest(h http.Handler) http.Handler
}

type Payload struct {
	TokenID   uuid.UUID `json:"token_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func newPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		TokenID:   tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

type BearerAuthorizer struct {
	secretKey []byte
}

func New(secretKey string) (Authorizer, error) {
	return &BearerAuthorizer{[]byte(secretKey)}, nil
}

func (b *BearerAuthorizer) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := newPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString(b.secretKey)
}

func (b *BearerAuthorizer) VerifyToken(token string) (*Payload, error) {
	const operation = ""
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return b.secretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		return nil, domain.Error(operation, err)
	}

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (a *BearerAuthorizer) AuthorizeRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// use responses package maybe?
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing auth header"))
			return
		}

		splitedAuthHeader := strings.Split(authHeader, " ")
		if len(splitedAuthHeader) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid auth header"))
			return
		}

		token := splitedAuthHeader[1]

		_, err := a.VerifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), accounts.TokenStr, token)
		instrumentation.Logger().WithField("TOKEN", ctx.Value(accounts.TokenStr)).Info("aaa")
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
