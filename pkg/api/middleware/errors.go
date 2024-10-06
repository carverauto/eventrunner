package middleware

import "errors"

var (
	errFailedToRetrieveCustomContext = errors.New("failed to retrieve custom context")
	errNoHandlerProvided             = errors.New("no handler provided")
)
