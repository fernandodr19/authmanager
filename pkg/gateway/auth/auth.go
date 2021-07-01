package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fernandodr19/library/pkg/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

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

type JWTAuthorizer struct {
	secretKey string
}

func New(secretKey string) (Authorizer, error) {
	return &JWTAuthorizer{secretKey}, nil
}

func (a *JWTAuthorizer) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := newPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(a.secretKey))
}

func (a *JWTAuthorizer) VerifyToken(token string) (*Payload, error) {
	const operation = ""
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(a.secretKey), nil
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

func (a *JWTAuthorizer) AuthorizeRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Let secure process the request. If it returns an error,
		// that indicates the request should not continue.
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			w.Write([]byte("Invalid auth header"))
			w.WriteHeader(http.StatusBadRequest)
		}
		token := authHeader[1]

		payload, err := a.VerifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Println(payload) //inject on context

		h.ServeHTTP(w, r)
	})
}
