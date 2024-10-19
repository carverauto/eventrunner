// Package middleware pkg/api/middleware/http_headers.go
package middleware

import (
	"context"
	"net/http"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	gofr "gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"
)

type contextKey string

const CustomContextKey contextKey = "customCtx"

func CustomHeadersMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use the existing context to create a custom context
			gofrCtx := &gofr.Context{
				Request: gofrHTTP.NewRequest(r),
			}
			customCtx := customctx.NewCustomContext(gofrCtx)

			// Extract headers from the HTTP request and store them in the custom context
			for key, values := range r.Header {
				if len(values) > 0 {
					customCtx.SetClaim(key, values[0])
				}
			}

			// Create a new context with the custom context
			ctxWithCustom := context.WithValue(r.Context(), CustomContextKey, customCtx)

			// Store the custom context back into the request
			r = r.WithContext(ctxWithCustom)

			// Call the next handler in the chain
			inner.ServeHTTP(w, r)
		})
	}
}
