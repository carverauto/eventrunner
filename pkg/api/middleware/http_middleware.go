package middleware

import (
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"gofr.dev/pkg/gofr"
)

// CombineMiddleware chains multiple middleware functions together.
func CombineMiddleware(middlewares ...interface{}) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		// Retrieve the custom context from the original context
		customCtx, ok := c.Request.Context().Value("customCtx").(*customctx.CustomContext)
		if !ok {
			return nil, errFailedToRetrieveCustomContext
		}

		// Define the final handler that will be called after applying all middleware
		finalHandler := func(customctx.Context) (interface{}, error) {
			return nil, errNoHandlerProvided
		}

		// Apply middlewares in reverse order to build the middleware chain
		for i := len(middlewares) - 1; i >= 0; i-- {
			switch m := middlewares[i].(type) {
			case func(customctx.Context) (interface{}, error):
				// Set the final handler to the current one if no other handler is set
				if i == len(middlewares)-1 {
					finalHandler = m
				} else {
					// Wrap the final handler in the current function
					finalHandler = func(ctx customctx.Context) (interface{}, error) {
						return m(ctx)
					}
				}
			case func(func(customctx.Context) (interface{}, error)) func(customctx.Context) (interface{}, error):
				// Wrap the final handler in middleware if it's a middleware function
				finalHandler = m(finalHandler)
			}
		}

		// Execute the final middleware chain with the custom context
		return finalHandler(customCtx)
	}
}

// Adapt converts a handler func and middlewares into a gofr.Handler.
func Adapt(h interface{}, middlewares ...handlers.Middleware) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		var handler handlers.Handler
		switch h := h.(type) {
		case func(*gofr.Context) (interface{}, error):
			handler = handlers.HandlerFunc(h)
		default:
			if h, ok := h.(handlers.Handler); ok {
				handler = h
			} else {
				panic("unsupported handler type")
			}
		}

		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}

		return handler.Handle(c)
	}
}
