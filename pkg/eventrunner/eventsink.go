package eventrunner

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type LogEventSink struct {
	// Add any necessary fields (e.g., database connection)
}

func (*LogEventSink) LogEvent(context.Context, *cloudevents.Event) error {
	// Implement event logging logic here
	return nil
}
