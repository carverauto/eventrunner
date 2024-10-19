package middleware

import (
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"gofr.dev/pkg/gofr"
)

// Adapt converts a handler func and middlewares into a gofr.Handler
func Adapt(h handlers.HandlerFunc, middlewares ...handlers.Middleware) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		handler := handlers.Handler(h)
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler.Handle(c)
	}
}
