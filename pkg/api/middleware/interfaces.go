package middleware

import "github.com/google/uuid"

//go:generate mockgen -destination=mock_custom_context.go -package=middleware github.com/carverauto/eventrunner/pkg/api/middleware CustomContext

type CustomContext interface {
	GetAPIKey() (string, bool)
	FindAPIKey(apiKey string) (uuid.UUID, uuid.UUID, error)
	SetClaim(key string, value interface{})
}
