package middleware

import (
	"context"

	"net/http"

	gofrHTTP "gofr.dev/pkg/gofr/http"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// APIKeyContextKey is the key used to store the API key in the context.
const APIKeyContextKey contextKey = "APIKey"

// APIKeyMiddleware is a middleware that extracts the API key from the request headers and stores it in the context.
func APIKeyMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			ctx := context.WithValue(r.Context(), APIKeyContextKey, apiKey)
			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAPIKeyFromContext retrieves the API key from the context.
func GetAPIKeyFromContext(ctx context.Context) (string, bool) {
	apiKey, ok := ctx.Value(APIKeyContextKey).(string)

	return apiKey, ok
}
