package middleware

import (
	"github.com/carverauto/eventrunner/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	gofrMiddleware "gofr.dev/pkg/gofr/http/middleware"
)

// Handler is an interface that wraps the basic Handle method.
type Handler interface {
	Handle(*gofr.Context) (interface{}, error)
}

// HandlerFunc is a function type that implements the Handler interface.
type HandlerFunc func(*gofr.Context) (interface{}, error)

// Handle calls f(c).
func (f HandlerFunc) Handle(c *gofr.Context) (interface{}, error) {
	return f(c)
}

// Middleware defines the standard middleware signature.
type Middleware func(Handler) Handler

// Adapt converts a HandlerFunc and middlewares into a gofr.Handler
func Adapt(h HandlerFunc, middlewares ...Middleware) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		var handler Handler = h
		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}
		return handler.Handle(c)
	}
}

// RequireRole is a middleware that checks if the user has the required role.
func RequireRole(roles ...string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(c *gofr.Context) (interface{}, error) {
			// Retrieve the JWT claim from the context
			claimData := c.Context.Value(gofrMiddleware.JWTClaim)

			// Assert that the claimData is of type jwt.MapClaims
			claims, ok := claimData.(jwt.MapClaims)
			if !ok {
				return nil, errors.NewInvalidParamError("JWT claims")
			}

			// Check if the user has the required role
			userRole, ok := claims["role"].(string)
			if !ok {
				return nil, errors.NewMissingParamError("role")
			}

			for _, role := range roles {
				if userRole == role {
					return next.Handle(c)
				}
			}

			return nil, errors.NewForbiddenError("Insufficient permissions")
		})
	}
}

// GetJWTClaims is a helper function to retrieve JWT claims from the context.
func GetJWTClaims(c *gofr.Context) (jwt.MapClaims, error) {
	claimData := c.Context.Value(gofrMiddleware.JWTClaim)
	claims, ok := claimData.(jwt.MapClaims)
	if !ok {
		return nil, errors.NewInvalidParamError("JWT claims")
	}
	return claims, nil
}
