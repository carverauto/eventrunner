package middleware

import (
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
	"gofr.dev/pkg/gofr"
)

func AuthenticateAPIKey(next handlers.Handler) handlers.Handler {
	return handlers.HandlerFunc(func(c *gofr.Context) (interface{}, error) {
		return authenticateAPIKey(customctx.NewCustomContext(c), next)
	})
}

func authenticateAPIKey(cc customctx.Context, next handlers.Handler) (interface{}, error) {
	apiKey, ok := cc.GetAPIKey()
	if !ok || apiKey == "" {
		return nil, eventingest.NewAuthError("Missing API Key")
	}

	tenantID, ok := cc.GetUUIDClaim("tenant_id")
	if !ok {
		return nil, eventingest.NewAuthError("Missing tenant ID")
	}

	customerID, ok := cc.GetUUIDClaim("customer_id")
	if !ok {
		return nil, eventingest.NewAuthError("Missing customer ID")
	}

	cc.SetClaim("api_key", apiKey)
	cc.SetClaim("tenant_id", tenantID)
	cc.SetClaim("customer_id", customerID)

	return next.Handle(cc.Context())
}

// RequireRole checks if the user has the required role to access the resource, otherwise returns an error.
func RequireRole(roles ...string) handlers.Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(c *gofr.Context) (interface{}, error) {
			cc := customctx.NewCustomContext(c)

			userRole, ok := cc.GetStringClaim("user_role")
			if !ok {
				return nil, eventingest.NewAuthError("Missing user role")
			}

			// Check if the user has the required role
			for _, role := range roles {
				if userRole == role {
					return next.Handle(c)
				}
			}

			return nil, eventingest.NewAuthError("Insufficient permissions")
		})
	}
}
