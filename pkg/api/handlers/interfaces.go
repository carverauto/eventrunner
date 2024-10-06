package handlers

import "gofr.dev/pkg/gofr"

//go:generate mockgen -destination=mock_handlers.go -package=handlers github.com/carverauto/eventrunner/pkg/api/handlers Handler

// Handler is an interface that wraps the basic Handle method.
type Handler interface {
	Handle(*gofr.Context) (interface{}, error)
}
