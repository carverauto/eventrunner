// Package eventrunner pkg/eventrunner/interfaces.go
package eventrunner

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
)

//go:generate mockgen -destination=mock_eventrunner.go -package=eventrunner -source=./interfaces.go EventSink,Consumer

type EventSink interface {
	LogEvent(context.Context, *cloudevents.Event) error
}

type Consumer interface {
	ConsumeEvent(ctx *gofr.Context, event *cloudevents.Event) error
}
