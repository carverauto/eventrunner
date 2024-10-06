package middleware

import (
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
)

// AuthenticateAPIKey checks if the API key is valid and active, otherwise returns an error.
func AuthenticateAPIKey(next func(CustomContext) (interface{}, error)) func(CustomContext) (interface{}, error) {
	return func(cc CustomContext) (interface{}, error) {
		apiKey, ok := cc.GetAPIKey()
		if !ok || apiKey == "" {
			return nil, eventingest.NewAuthError("Missing API Key")
		}

		tenantID, customerID, err := cc.FindAPIKey(apiKey)
		if err != nil {
			return nil, eventingest.NewAuthError("Invalid API Key")
		}

		cc.SetClaim("api_key", apiKey)
		cc.SetClaim("tenant_id", tenantID)
		cc.SetClaim("customer_id", customerID)

		return next(cc)
	}
}

// RequireRole checks if the user has the required role to access the resource, otherwise returns an error.
// The user's role is stored in the JWT token. The roles parameter is a list of roles that are allowed
// to access the resource.
func RequireRole(roles ...string) func(
	func(customctx.Context) (interface{}, error)) func(customctx.Context) (interface{}, error) {
	return func(
		next func(customctx.Context) (interface{}, error)) func(customctx.Context) (interface{}, error) {
		return func(cc customctx.Context) (interface{}, error) {
			userRole, ok := cc.GetStringClaim("user_role")
			if !ok {
				return nil, eventingest.NewAuthError("Missing user role")
			}

			// Check if the user has the required role
			for _, role := range roles {
				if userRole == role {
					return next(cc)
				}
			}

			return nil, eventingest.NewAuthError("Insufficient permissions")
		}
	}
}
