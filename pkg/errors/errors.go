package errors

import (
	"fmt"
	"net/http"
)

// AppError is the base error type for our application
type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code for the error
func (e AppError) StatusCode() int {
	return e.Code
}

// NewAppError creates a new AppError
func NewAppError(code int, message string) AppError {
	return AppError{Code: code, Message: message}
}

// InvalidParamError represents an error due to invalid parameters
type InvalidParamError struct {
	AppError
	Param []string
}

// NewInvalidParamError creates a new InvalidParamError
func NewInvalidParamError(params ...string) InvalidParamError {
	return InvalidParamError{
		AppError: AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid parameter(s): %v", params),
		},
		Param: params,
	}
}

// MissingParamError represents an error due to missing parameters
type MissingParamError struct {
	AppError
	Param []string
}

// NewMissingParamError creates a new MissingParamError
func NewMissingParamError(params ...string) MissingParamError {
	return MissingParamError{
		AppError: AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Missing parameter(s): %v", params),
		},
		Param: params,
	}
}

// DatabaseError represents a database-related error
type DatabaseError struct {
	AppError
	Err error
}

// NewDatabaseError creates a new DatabaseError
func NewDatabaseError(err error, message string) DatabaseError {
	return DatabaseError{
		AppError: AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Database error: %s", message),
		},
		Err: err,
	}
}

// UnauthorizedError represents an unauthorized access error
type UnauthorizedError struct {
	AppError
}

// NewUnauthorizedError creates a new UnauthorizedError
func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		AppError: AppError{
			Code:    http.StatusUnauthorized,
			Message: message,
		},
	}
}

// ForbiddenError represents a forbidden access error
type ForbiddenError struct {
	AppError
}

// NewForbiddenError creates a new ForbiddenError
func NewForbiddenError(message string) ForbiddenError {
	return ForbiddenError{
		AppError: AppError{
			Code:    http.StatusForbidden,
			Message: message,
		},
	}
}
