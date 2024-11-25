package middleware

import (
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/errors"
	"gofr.dev/pkg/gofr"
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

// getCustomContext extracts the custom context from the request context
func getCustomContext(c *gofr.Context) (customctx.Context, error) {
	ctx := c.Request.Context()
	customCtxVal := ctx.Value("customCtx")
	if customCtxVal == nil {
		return nil, errors.NewAppError(500, "custom context not found")
	}

	customCtx, ok := customCtxVal.(customctx.Context)
	if !ok {
		return nil, errors.NewAppError(500, "invalid custom context type")
	}

	return customCtx, nil
}

// RequireRole is a middleware that checks if the user has any of the required roles
func RequireRole(roles ...string) Middleware {
	return func(next Handler) Handler {
		return HandlerFunc(func(c *gofr.Context) (interface{}, error) {
			customCtx, err := getCustomContext(c)
			if err != nil {
				return nil, err
			}

			// Get role from X-User-Role header (stored as claim)
			userRole, ok := customCtx.GetStringClaim("X-User-Role")
			if !ok {
				return nil, errors.NewMissingParamError("X-User-Role header")
			}

			// Check if user has any of the required roles
			hasRole := false
			for _, role := range roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return nil, errors.NewForbiddenError("Insufficient permissions")
			}

			return next.Handle(c)
		})
	}
}

// RequireTenant is a middleware that ensures a tenant ID is present
func RequireTenant(next Handler) Handler {
	return HandlerFunc(func(c *gofr.Context) (interface{}, error) {
		customCtx, err := getCustomContext(c)
		if err != nil {
			return nil, err
		}

		tenantID, ok := customCtx.GetStringClaim("X-Tenant-ID")
		if !ok || tenantID == "" {
			return nil, errors.NewMissingParamError("X-Tenant-ID header")
		}

		return next.Handle(c)
	})
}

// RequireUser is a middleware that ensures a user ID is present
func RequireUser(next Handler) Handler {
	return HandlerFunc(func(c *gofr.Context) (interface{}, error) {
		customCtx, err := getCustomContext(c)
		if err != nil {
			return nil, err
		}

		userID, ok := customCtx.GetStringClaim("X-User-Id")
		if !ok || userID == "" {
			return nil, errors.NewMissingParamError("X-User-Id header")
		}

		return next.Handle(c)
	})
}

// GetAuthClaims extracts the relevant authentication claims from the custom context
func GetAuthClaims(c *gofr.Context) (map[string]string, error) {
	customCtx, err := getCustomContext(c)
	if err != nil {
		return nil, err
	}

	claims := make(map[string]string)

	// Extract common auth claims
	if userID, ok := customCtx.GetStringClaim("X-User-Id"); ok {
		claims["user_id"] = userID
	}
	if tenantID, ok := customCtx.GetStringClaim("X-Tenant-Id"); ok {
		claims["tenant_id"] = tenantID
	}
	if userRole, ok := customCtx.GetStringClaim("X-User-Role"); ok {
		claims["user_role"] = userRole
	}
	if email, ok := customCtx.GetStringClaim("X-User-Email"); ok {
		claims["email"] = email
	}

	return claims, nil
}
