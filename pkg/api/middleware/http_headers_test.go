package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestCustomHeadersMiddleware(t *testing.T) {
	// Create a sample HTTP request with headers
	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	req.Header.Set("X-Custom-Header", "CustomValue")
	req.Header.Set("Authorization", "Bearer token")

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a mock final handler that will be called after the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the custom context from the request context
		customCtx, ok := r.Context().Value(CustomContextKey).(*customctx.CustomContext)
		if !ok {
			t.Errorf("Failed to retrieve custom context")
			return
		}

		// Validate that the headers were correctly set in the custom context
		headerValue, ok := customCtx.GetClaim("X-Custom-Header")
		assert.True(t, ok)
		assert.Equal(t, "CustomValue", headerValue)

		authValue, ok := customCtx.GetClaim("Authorization")
		assert.True(t, ok)
		assert.Equal(t, "Bearer token", authValue)

		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the CustomHeadersMiddleware
	middleware := CustomHeadersMiddleware()
	wrappedHandler := middleware(handler)

	// Serve the request
	wrappedHandler.ServeHTTP(rr, req)

	// Verify the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}
