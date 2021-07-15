package authorizer

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fernandodr19/authmanager/pkg/domain"
	"github.com/fernandodr19/authmanager/pkg/domain/entities/accounts"
	acc_usecase "github.com/fernandodr19/authmanager/pkg/domain/usecases/accounts"
	"github.com/fernandodr19/authmanager/pkg/domain/vos"
	"github.com/fernandodr19/authmanager/pkg/gateway/api/responses"
	"github.com/fernandodr19/authmanager/pkg/instrumentation/logger"
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
	TokenID   uuid.UUID `json:"token_id"`
	UserID    vos.AccID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// Valid checks if a payload is valid
func (p *payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// CreateTokens generate both access & refresh tokens
func (b *bearerAuthorizer) CreateTokens(acc accounts.Account, accessDuration time.Duration, refreshDuration time.Duration) (vos.Tokens, error) {
	const operation = "authorizer.BearerAuthorizer.CreateToken"

	userID := acc.ID

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

func (b *bearerAuthorizer) createToken(userID vos.AccID, duration time.Duration) (string, error) {
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
			logger.FromCtx(r.Context()).Infoln("missing token")
			resp := responses.Unauthorized(nil, responses.ErrMissingAuthHeader)
			responses.SendJSON(w, resp.Payload, http.StatusUnauthorized)
			return
		}

		splitedAuthHeader := strings.Split(authHeader, " ")
		if len(splitedAuthHeader) != 2 {
			logger.FromCtx(r.Context()).Infoln("invalid token")
			resp := responses.Unauthorized(nil, responses.ErrInvalidAuthHeader)
			responses.SendJSON(w, resp.Payload, http.StatusUnauthorized)
			return
		}

		token := splitedAuthHeader[1]

		payload, err := b.verifyToken(token)
		if err != nil {
			logger.FromCtx(r.Context()).WithError(err).Infoln("unauthorized")
			if errors.Is(err, ErrExpiredToken) {
				resp := responses.Unauthorized(nil, responses.ErrExpiredToken)
				responses.SendJSON(w, resp.Payload, http.StatusUnauthorized)
				return
			}
			resp := responses.Unauthorized(nil, responses.ErrUnauthorized)
			responses.SendJSON(w, resp.Payload, http.StatusUnauthorized)
			return
		}

		// TODO: check param user id against token's?

		ctx := context.WithValue(r.Context(), accounts.UserIDContextKey, payload.UserID)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (b *bearerAuthorizer) verifyToken(token string) (*payload, error) {
	const operation = "authorizer.bearerAuthorizer.verifyToken"

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, domain.Error(operation, ErrInvalidToken)
		}
		return b.secretKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &payload{}, keyFunc)
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
