package context

import (
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

//go:generate mockgen -destination=mock_context.go -package=context github.com/carverauto/eventrunner/pkg/context Interface

type Interface interface {
	SetClaim(key string, value interface{})
	GetClaim(key string) (interface{}, bool)
	GetStringClaim(key string) (string, bool)
	GetUUIDClaim(key string) (uuid.UUID, bool)
	GetAPIKey() (string, bool)
	Bind(v interface{}) error
	Context() *gofr.Context
}
