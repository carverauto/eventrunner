package eventrunner

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
)

type EventSink interface {
	LogEvent(context.Context, *cloudevents.Event) error
}

type Consumer interface {
	ConsumeEvent(ctx *gofr.Context, event *cloudevents.Event) error
}
