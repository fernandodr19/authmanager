package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fernandodr19/library/pkg/domain/entities/accounts"
	accounts_uc "github.com/fernandodr19/library/pkg/domain/usecases/accounts"
)

// Response represents an API response
type Response struct {
	Status  int
	Error   error
	Payload interface{}
	headers map[string]string
}

// Headers list response headers
func (r *Response) Headers() map[string]string {
	return r.headers
}

// SetHeader set response header
func (r *Response) SetHeader(key, value string) {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}

	r.headers[key] = value
}

type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ErrorPayload represents error response payload
type ErrorPayload struct {
	Error `json:"errors"`
}

// shared
var (
	ErrInternalServerError = ErrorPayload{Error: Error{Code: "error:internal_server_error", Description: "Internal Server Error"}}
	ErrInvalidBody         = ErrorPayload{Error: Error{Code: "error:invalid_body", Description: "Invalid body"}}
	ErrInvalidAuth         = ErrorPayload{Error: Error{Code: "error:invalid_auth", Description: "Invalid authorization"}}
	ErrInvalidParams       = ErrorPayload{Error: Error{Code: "error:invalid_parameters", Description: "Invalid query parameters"}}
	ErrNotImplemented      = ErrorPayload{Error: Error{Code: "error:not_implemented", Description: "Not implemented"}}
)

// accounts
var (
	ErrAccountNotFound        = ErrorPayload{Error: Error{Code: "error:account_not_found", Description: "Account not found"}}
	ErrInvalidEmail           = ErrorPayload{Error: Error{Code: "error:invalid_email", Description: "Invalid email"}}
	ErrInvalidPassword        = ErrorPayload{Error: Error{Code: "error:invalid_password", Description: "Invalid password"}}
	ErrWrongPassword          = ErrorPayload{Error: Error{Code: "error:wrong_password", Description: "Wrong password"}}
	ErrEmailAlreadyRegistered = ErrorPayload{Error: Error{Code: "error:already_registered", Description: "Email already registered"}}
)

// ErrorResponse maps response error
func ErrorResponse(err error) Response {
	switch {
	case errors.Is(err, accounts.ErrInvalidEmail):
		return UnprocessableEntity(err, ErrInvalidEmail)
	case errors.Is(err, accounts.ErrInvalidPassword):
		return UnprocessableEntity(err, ErrInvalidPassword)
	case errors.Is(err, accounts_uc.ErrNotImplemented):
		return NotImplemented(err)
	case errors.Is(err, accounts_uc.ErrEmailAlreadyRegistered):
		return Conflict(err, ErrEmailAlreadyRegistered)
	case errors.Is(err, accounts_uc.ErrAccountNotFound):
		return NotFound(err, ErrAccountNotFound)
	case errors.Is(err, accounts_uc.ErrWrongPassword):
		return Unauthorized(err, ErrWrongPassword)
	default:
		return InternalServerError(err)
	}
}

// InternalServerError 500
func InternalServerError(err error) Response {
	return Response{
		Status:  http.StatusInternalServerError,
		Error:   err,
		Payload: ErrInternalServerError,
	}
}

// NotImplemented 501
func NotImplemented(err error) Response {
	return Response{
		Status:  http.StatusNotImplemented,
		Error:   err,
		Payload: ErrNotImplemented,
	}
}

// BadRequest 400
func BadRequest(err error, payload ErrorPayload) Response {
	return genericError(http.StatusBadRequest, err, payload)
}

// Unauthorized 401
func Unauthorized(err error, payload ErrorPayload) Response {
	return genericError(http.StatusUnauthorized, err, payload)
}

// NotFound 404
func NotFound(err error, payload ErrorPayload) Response {
	return genericError(http.StatusNotFound, err, payload)
}

// Conflict 409
func Conflict(err error, payload ErrorPayload) Response {
	return genericError(http.StatusConflict, err, payload)
}

// UnprocessableEntity 422
func UnprocessableEntity(err error, payload ErrorPayload) Response {
	return genericError(http.StatusUnprocessableEntity, err, payload)
}

func genericError(status int, err error, payload ErrorPayload) Response {
	return Response{
		Status:  status,
		Error:   err,
		Payload: payload,
	}
}

// NoContent 204
func NoContent() Response {
	return Response{
		Status: http.StatusNoContent,
	}
}

// OK 200
func OK(payload interface{}) Response {
	return Response{
		Status:  http.StatusOK,
		Payload: payload,
	}
}

// Created 201
func Created(payload interface{}) Response {
	return Response{
		Status:  http.StatusCreated,
		Payload: payload,
	}
}

// Accepted 202
func Accepted(payload interface{}) Response {
	return Response{
		Status:  http.StatusAccepted,
		Payload: payload,
	}
}

// SendJSON responds requests based on
func SendJSON(w http.ResponseWriter, payload interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil { // Blank body is not valid JSON.
		return nil
	}

	switch p := payload.(type) {
	case string:
		if p == "" {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(payload)
}
