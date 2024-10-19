package middleware

import (
	"net/http"
	"strings"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	ory "github.com/ory/client-go"
	"gofr.dev/pkg/gofr"
)

func OryAuthMiddleware(client *ory.APIClient) func(next handlers.Handler) handlers.Handler {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(c *gofr.Context) (interface{}, error) {
			cc := customctx.NewCustomContext(c)
			token := cc.GetHeader("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")

			if token == "" {
				return nil, NewErrorResponse(http.StatusUnauthorized, "Missing authorization token")
			}

			introspect, _, err := client.OAuth2API.IntrospectOAuth2Token(c.Context).Token(token).Execute()
			if err != nil || !introspect.Active {
				return nil, NewErrorResponse(http.StatusUnauthorized, "Invalid or expired token")
			}

			cc.SetClaim("sub", introspect.Sub)
			cc.SetClaim("scope", introspect.Scope)

			return next.Handle(cc.Context())
		})
	}
}
