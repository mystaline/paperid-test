package helper

import (
	"github.com/gofiber/fiber/v3"
)

var _ error = (*Error)(nil)

// Error represents a structured HTTP error with status code and message.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// ErrInternal creates a 500 Internal Server Error.
func ErrInternal(message string) *Error {
	return &Error{
		Code:    fiber.StatusInternalServerError,
		Message: message,
	}
}

// ErrInternal creates a 422 Internal Server Error.
func ErrUnprocessableEntity(message string) *Error {
	return &Error{
		Code:    fiber.StatusUnprocessableEntity,
		Message: message,
	}
}

// ErrBadRequest creates a 400 Bad Request error.
func ErrBadRequest(message string) *Error {
	return &Error{
		Code:    fiber.StatusBadRequest,
		Message: message,
	}
}

// ErrNotFound creates a 404 Not Found error.
func ErrNotFound(message string) *Error {
	return &Error{
		Code:    fiber.StatusNotFound,
		Message: message,
	}
}
