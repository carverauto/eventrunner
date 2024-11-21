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
