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

// CassandraInsertError is a custom error type for Cassandra insertion errors.
type CassandraInsertError struct {
	OriginalError error
}

// Error implements the error interface for CassandraInsertError.
func (cie *CassandraInsertError) Error() string {
	return fmt.Sprintf("failed to insert event into Cassandra: %v", cie.OriginalError)
}

// Unwrap allows errors.Is and errors.As to work with CassandraInsertError.
func (cie *CassandraInsertError) Unwrap() error {
	return cie.OriginalError
}

func (*CassandraEventSink) ConsumeEvent(ctx *gofr.Context, event *cloudevents.Event) error {
	if ctx == nil {
		return errNilContext
	}

	if event == nil {
		return errNilEvent
	}

	if ctx.Cassandra == nil {
		return errNilCassandra
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
		return &CassandraInsertError{OriginalError: err}
	}

	return nil
}
