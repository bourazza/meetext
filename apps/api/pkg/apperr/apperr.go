package apperr

import (
	"errors"
	"net/http"
)

// AppError is a domain-level error with an HTTP status code and machine-readable code.
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func New(statusCode int, code, message string) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message}
}

func Wrap(err error, statusCode int, code, message string) *AppError {
	return &AppError{StatusCode: statusCode, Code: code, Message: message, Err: err}
}

var (
	ErrNotFound           = New(http.StatusNotFound, "NOT_FOUND", "resource not found")
	ErrUnauthorized       = New(http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	ErrForbidden          = New(http.StatusForbidden, "FORBIDDEN", "forbidden")
	ErrConflict           = New(http.StatusConflict, "CONFLICT", "resource already exists")
	ErrBadRequest         = New(http.StatusBadRequest, "BAD_REQUEST", "bad request")
	ErrUnprocessable      = New(http.StatusUnprocessableEntity, "UNPROCESSABLE", "unprocessable entity")
	ErrInternal           = New(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
	ErrEmailNotVerified   = New(http.StatusForbidden, "EMAIL_NOT_VERIFIED", "verify your email before signing in")
	ErrTokenExpired       = New(http.StatusUnauthorized, "TOKEN_EXPIRED", "token has expired")
	ErrTokenInvalid       = New(http.StatusUnauthorized, "TOKEN_INVALID", "token is invalid")
	ErrFileTooLarge       = New(http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "file exceeds maximum allowed size")
	ErrUnsupportedFile    = New(http.StatusBadRequest, "UNSUPPORTED_FILE", "unsupported file type")
)

func As(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
