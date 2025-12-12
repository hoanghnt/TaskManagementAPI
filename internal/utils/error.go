package utils

import "errors"

// Custom error types
var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrBadRequest         = errors.New("bad request")
	ErrInternalServer     = errors.New("internal server error")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Task & Category specific errors
	ErrCategoryNotFound = errors.New("category not found")
	ErrTaskNotFound     = errors.New("task not found")
)

// IsNotFoundError checks if error is not found error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsDuplicateError checks if error is duplicate error
func IsDuplicateError(err error) bool {
	return errors.Is(err, ErrDuplicateEntry)
}
