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
	verifier IDTokenVerifier
	config   *config.OAuthConfig
}

// NewJWTMiddleware creates a new JWTMiddleware.
func NewJWTMiddleware(ctx context.Context, cfg *config.OAuthConfig) (*JWTMiddleware, error) {
	provider, err := oidc.NewProvider(ctx, cfg.KeycloakURL)
	if err != nil {
		return nil, eventingest.NewInternalError("Failed to get provider")
	}

	oidcVerifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})
	verifier := NewOIDCVerifier(oidcVerifier)

	return &JWTMiddleware{
		verifier: verifier,
		config:   cfg,
	}, nil
}

// Validate validates a JWT token.
func (m *JWTMiddleware) Validate(next func(customctx.Context) (interface{}, error)) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		cc := customctx.NewCustomContext(c)

		// Safely retrieve Authorization header from context
		authHeaderValue := c.Request.Context().Value("Authorization")

		authHeader, ok := authHeaderValue.(string)
		if !ok || authHeader == "" {
			return nil, eventingest.NewAuthError("Missing or invalid authorization header")
		}

		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")
		if bearerToken == authHeader {
			return nil, eventingest.NewAuthError("Invalid authorization header format")
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
