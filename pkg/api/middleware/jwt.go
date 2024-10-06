package middleware

import (
	"context"
	"strings"

	"github.com/carverauto/eventrunner/pkg/config"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
	"github.com/coreos/go-oidc/v3/oidc"
	"gofr.dev/pkg/gofr"
)

// JWTMiddleware is a middleware that validates JWT tokens.
type JWTMiddleware struct {
	verifier *oidc.IDTokenVerifier
	config   *config.OAuthConfig
}

// NewJWTMiddleware creates a new JWTMiddleware.
func NewJWTMiddleware(ctx context.Context, config *config.OAuthConfig) (*JWTMiddleware, error) {
	provider, err := oidc.NewProvider(ctx, config.KeycloakURL)
	if err != nil {
		return nil, eventingest.NewInternalError("Failed to get provider")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	return &JWTMiddleware{
		verifier: verifier,
		config:   config,
	}, nil
}

// Validate validates a JWT token.
func (m *JWTMiddleware) Validate(next func(*customctx.Context) (interface{}, error)) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		cc := customctx.NewCustomContext(c)

		authHeader := c.Request.Context().Value("Authorization").(string)
		if authHeader == "" {
			return nil, eventingest.NewAuthError("Missing authorization header")
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == authHeader {
			return nil, eventingest.NewAuthError("Invalid authorization header")
		}

		token, err := m.verifier.Verify(context.Background(), bearerToken)
		if err != nil {
			return nil, eventingest.NewAuthError("Invalid token")
		}

		var claims struct {
			TenantID   string `json:"tenant_id"`
			CustomerID string `json:"customer_id"`
		}

		if err := token.Claims(&claims); err != nil {
			return nil, eventingest.NewInternalError("Failed to parse claims")
		}

		if claims.TenantID == "" || claims.CustomerID == "" {
			return nil, eventingest.NewAuthError("Missing required claims")
		}

		cc.SetClaim("tenant_id", claims.TenantID)
		cc.SetClaim("customer_id", claims.CustomerID)

		return next(cc)
	}
}
