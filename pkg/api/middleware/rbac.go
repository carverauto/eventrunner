package middleware

import (
	"github.com/carverauto/eventrunner/pkg/api"
	"github.com/carverauto/eventrunner/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"gofr.dev/pkg/gofr"
	gofrHTTP "gofr.dev/pkg/gofr/http"
)

func AuthenticateAPIKey(next gofr.Handler) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		apiKey := c.Request.Param("X-API-Key")
		if apiKey == "" {
			return nil, gofrHTTP.ErrorMissingParam{Params: []string{"X-API-Key"}}
		}

		var key models.APIKey
		if err := c.Mongo.FindOne(c, "api_keys", bson.M{"key": apiKey, "active": true}, &key); err != nil {
			return nil, gofrHTTP.ErrorInvalidParam{Params: []string{"X-API-Key"}}
		}

		// Store the API key in the context for later use if needed
		// Note: There doesn't seem to be a direct method to set parameters in the context
		// You might need to implement a custom method or use a different approach to store this

		return next(c)
	}
}

func RequireRole(roles ...string) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		userRole := c.Request.Param("X-User-Role")
		if userRole == "" {
			return nil, gofrHTTP.ErrorMissingParam{Params: []string{"X-User-Role"}}
		}

		for _, role := range roles {
			if userRole == role {
				return nil, nil // Allow the request to proceed
			}
		}

		return nil, api.ErrInsufficientPermissions
	}
}
