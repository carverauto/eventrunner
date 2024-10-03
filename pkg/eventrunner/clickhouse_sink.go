// Package eventrunner pkg/eventrunner/clickhouse_sink.go
package eventrunner

import (
	"encoding/json"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
)

type ClickhouseEventSink struct {
	app *gofr.App
}

func NewClickhouseEventSink(app *gofr.App) *ClickhouseEventSink {
	return &ClickhouseEventSink{app: app}
}

func (s *ClickhouseEventSink) ConsumeEvent(ctx *gofr.Context, event *cloudevents.Event) error {
	// Convert the event to JSON
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Insert the event into Clickhouse
	query := `
        INSERT INTO events (id, source, type, subject, time, data_contenttype, data, specversion)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
	err = ctx.Clickhouse.Exec(ctx, query,
		event.ID(),
		event.Source(),
		event.Type(),
		event.Subject(),
		event.Time().Format(time.RFC3339),
		event.DataContentType(),
		string(eventJSON),
		event.SpecVersion(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert event into Clickhouse: %w", err)
	}

	return nil
}
