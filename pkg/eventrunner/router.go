// Package eventrunner pkg/eventrunner/router.go
package eventrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	nats "github.com/carverauto/gofr-nats"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type EventRouter struct {
	app        *gofr.App
	natsClient *nats.PubSubWrapper
	bufferPool *sync.Pool
}

func NewEventRouter() *EventRouter {
	app := gofr.New()

	subjects := strings.Split(os.Getenv("NATS_SUBJECTS"), ",")

	natsClient := nats.New(&nats.Config{
		Server: os.Getenv("PUBSUB_BROKER"),
		Stream: nats.StreamConfig{
			Stream:   os.Getenv("NATS_STREAM"),
			Subjects: subjects,
		},
		MaxWait:     5 * time.Second,
		BatchSize:   100,
		MaxPullWait: 10,
		Consumer:    os.Getenv("NATS_CONSUMER"),
		CredsFile:   os.Getenv("NATS_CREDS_FILE"),
	})
	natsClient.UseLogger(app.Logger)
	natsClient.UseMetrics(app.Metrics())
	natsClient.Connect()

	app.AddPubSub(natsClient)

	return &EventRouter{
		app:        app,
		natsClient: natsClient,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (er *EventRouter) Start() {
	er.app.Subscribe("eventrunner.>", er.handleEvent)
	er.app.Run()
}

func (er *EventRouter) handleEvent(c *gofr.Context) error {
	event := cloudevents.NewEvent()
	if err := c.Bind(&event); err != nil {
		return fmt.Errorf("failed to bind cloud event: %w", err)
	}

	eventType := event.Type()
	tenantID, err := event.Context.GetExtension("tenantid")
	if err != nil {
		return fmt.Errorf("missing tenant_id in event: %w", err)
	}

	// Route to appropriate consumer queue based on event type
	consumerQueue := fmt.Sprintf("events.%s", eventType)

	buf := er.getBuffer()
	defer er.putBuffer(buf)

	if err := json.NewEncoder(buf).Encode(event); err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	messageID := uuid.New().String()
	if err := er.natsClient.Publish(c.Context, consumerQueue, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	c.Logger.Info("Event routed successfully",
		"event_id", event.ID(),
		"message_id", messageID,
		"event_type", eventType,
		"tenant_id", tenantID,
		"consumer_queue", consumerQueue)

	return nil
}
