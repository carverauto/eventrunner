package middleware

import (
	"context"
	"net/http"

	gofrHTTP "gofr.dev/pkg/gofr/http"
)

func APIKeyMiddleware() gofrHTTP.Middleware {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			ctx := context.WithValue(r.Context(), "APIKey", apiKey)
			inner.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
