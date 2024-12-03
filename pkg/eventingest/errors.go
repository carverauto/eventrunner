/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

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
