package eventrunner

import (
	"encoding/json"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
)

type CassandraEventSink struct{}

func NewCassandraEventSink() *CassandraEventSink {
	return &CassandraEventSink{}
}

func (s *CassandraEventSink) ConsumeEvent(ctx *gofr.Context, event *cloudevents.Event) error {
	if ctx == nil {
		return fmt.Errorf("nil context provided to CassandraEventSink")
	}
	if event == nil {
		return fmt.Errorf("nil event provided to CassandraEventSink")
	}
	if ctx.Cassandra == nil {
		return fmt.Errorf("cassandra client is nil in CassandraEventSink")
	}

	// Get the event data as []byte
	dataJSON := event.Data()

	// If you need to ensure it's valid JSON, you can optionally validate it:
	var jsonObj interface{}
	if err := json.Unmarshal(dataJSON, &jsonObj); err != nil {
		return fmt.Errorf("event data is not valid JSON: %w", err)
	}

	// Insert the event into Cassandra
	query := `INSERT INTO events (id, source, type, subject, time, data_contenttype, data, specversion)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// Use the actual timestamp instead of a formatted string
	eventTime := event.Time()
	if eventTime.IsZero() {
		eventTime = time.Now() // Use current time if event time is not set
	}

	err := ctx.Cassandra.Exec(query,
		event.ID(),
		event.Source(),
		event.Type(),
		event.Subject(),
		eventTime,
		event.DataContentType(),
		string(dataJSON),
		event.SpecVersion(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert event into Cassandra: %w", err)
	}

	return nil
}
