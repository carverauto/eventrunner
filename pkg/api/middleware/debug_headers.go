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

// Package middleware pkg/api/middleware/debug_headers.go

package middleware

import (
	"net/http"
	"strings"

	gofrHTTP "gofr.dev/pkg/gofr/http"
	"gofr.dev/pkg/gofr/logging"
)

func DebugHeadersMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := logging.NewLogger(logging.DEBUG)

			logger.Debug("-------------------- Request Headers --------------------")

			// Log ALL headers
			for name, values := range r.Header {
				for _, value := range values {
					if strings.ToLower(name) == "authorization" {
						// Truncate auth header for security
						logger.Debugf("%s: %s...", name, value[:debugMin(len(value), 15)])
					} else {
						logger.Debugf("%s: %s", name, value)
					}
				}
			}

			logger.Debug("-----------------------------------------------------")

			inner.ServeHTTP(w, r)
		})
	}
}

/*
func DebugHeadersMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We'll use stdout logger for debugging
			logger := logging.NewLogger(logging.DEBUG)

			logger.Debug("-------------------- Request Headers --------------------")

			// Log important headers we're particularly interested in
			importantHeaders := []string{
				"X-Tenant-ID",
				"X-User-ID",
				"X-User-Email",
				"X-Correlation-ID",
				"Authorization",
			}

			for _, header := range importantHeaders {
				if value := r.Header.Get(header); value != "" {
					if header == "Authorization" {
						// Truncate auth header for security
						logger.Debugf("%s: %s...", header, value[:debugMin(len(value), 15)])
					} else {
						logger.Debugf("%s: %s", header, value)
					}
				} else {
					logger.Debugf("%s: not present", header)
				}
			}

			logger.Debug("-----------------------------------------------------")

			// Call the next handler in the chain
			inner.ServeHTTP(w, r)
		})
	}
}

*/

func debugMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
