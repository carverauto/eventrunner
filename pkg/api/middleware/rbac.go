package middleware

import (
	"fmt"
	"net/http"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	gofrMiddleware "gofr.dev/pkg/gofr/http/middleware"
)

func RequireRole(roles ...string) func(next handlers.Handler) handlers.Handler {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(c *gofr.Context) (interface{}, error) {
			// Retrieve JWT claims
			claimData := c.Context.Value(gofrMiddleware.JWTClaim)
			claims, ok := claimData.(jwt.MapClaims)
			if !ok {
				return nil, fmt.Errorf("invalid claim data type")
			}

			// Check if the user has the required role
			userRoles, ok := claims["roles"].([]interface{})
			if !ok {
				return nil, NewErrorResponse(http.StatusForbidden, "Missing roles claim")
			}

			for _, role := range roles {
				for _, userRole := range userRoles {
					if userRole.(string) == role {
						return next.Handle(c)
					}
				}
			}

			return nil, NewErrorResponse(http.StatusForbidden, "Insufficient permissions")
		})
	}
}

// GetJWTClaims helper function
func GetJWTClaims(c *gofr.Context) (jwt.MapClaims, error) {
	claimData := c.Context.Value(gofrMiddleware.JWTClaim)
	claims, ok := claimData.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claim data type")
	}
	return claims, nil
}
