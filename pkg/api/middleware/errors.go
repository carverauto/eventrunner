package middleware

import "errors"

var (
	errFailedToRetrieveCustomContext = errors.New("failed to retrieve custom context")
	errNoHandlerProvided             = errors.New("no handler provided")
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (j ErrorResponse) Error() string {
	return j.Message
}

func NewErrorResponse(status int, message string) error {
	return ErrorResponse{Status: status, Message: message}
}
