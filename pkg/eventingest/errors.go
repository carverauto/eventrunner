package eventingest

import (
	"fmt"
	"net/http"
)

var (
	errInvalidJSON = NewProcessingError("Invalid JSON")
	errForwardFail = NewProcessingError("Failed to forward event")
)

// CustomError is a custom error type that includes an HTTP status code.
type CustomError struct {
	Message    string
	StatusCode int
}

// Error returns the error message.
func (e CustomError) Error() string {
	return e.Message
}

// NewAuthError creates a new CustomError with a 401 status code.
func NewAuthError(message string) CustomError {
	return CustomError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewInternalError creates a new CustomError with a 500 status code.
func NewInternalError(message string) CustomError {
	return CustomError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewProcessingError creates a new CustomError with a 500 status code.
func NewProcessingError(message string) CustomError {
	return CustomError{
		Message:    fmt.Sprintf("Failed to process event: %s", message),
		StatusCode: http.StatusInternalServerError,
	}
}
