package middleware

import (
	"context"

	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

//go:generate mockgen -destination=mock_middleware.go -package=middleware github.com/carverauto/eventrunner/pkg/api/middleware CustomContext,JWTMiddlewareInterface,IDTokenVerifier,Token

// CustomContext interface for mocking CustomContext.
type CustomContext interface {
	GetAPIKey() (string, bool)
	FindAPIKey(apiKey string) (uuid.UUID, uuid.UUID, error)
	SetClaim(key string, value interface{})
}

// JWTMiddlewareInterface interface for mocking JWTMiddleware.
type JWTMiddlewareInterface interface {
	Validate(next func(customctx.Context) (interface{}, error)) gofr.Handler
}

// IDTokenVerifier interface for mocking oidc.IDTokenVerifier.
type IDTokenVerifier interface {
	Verify(ctx context.Context, rawToken string) (Token, error)
}

// Token interface to abstract oidc.IDToken.
type Token interface {
	Claims(v interface{}) error
}
