package middleware

import (
	"net/http"
	"strings"

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
func RequireRole(roles ...string) func(next handlers.Handler) handlers.Handler {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(c *gofr.Context) (interface{}, error) {
			cc := customctx.NewCustomContext(c)
			scope, ok := cc.GetStringClaim("scope")
			if !ok {
				return nil, NewErrorResponse(http.StatusForbidden, "Missing scope claim")
			}

			scopes := strings.Split(scope, " ")
			for _, role := range roles {
				for _, s := range scopes {
					if s == role {
						return next.Handle(c)
					}
				}
			}

			return nil, NewErrorResponse(http.StatusForbidden, "Insufficient permissions")
		})
	}
}
