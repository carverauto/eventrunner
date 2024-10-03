// Package eventrunner pkg/eventrunner/cassandra_sink.go
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
	// Convert the event data to JSON
	dataJSON, err := json.Marshal(event.Data())
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Insert the event into Cassandra
	query := `INSERT INTO events (id, source, type, subject, time, data_contenttype, data, specversion)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	err = ctx.Cassandra.Exec(query,
		event.ID(),
		event.Source(),
		event.Type(),
		event.Subject(),
		event.Time().Format(time.RFC3339),
		event.DataContentType(),
		string(dataJSON),
		event.SpecVersion(),
	)
	if err != nil {
		return fmt.Errorf("failed to insert event into Cassandra: %w", err)
	}

	return nil
}
