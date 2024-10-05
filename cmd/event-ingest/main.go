package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/carverauto/gofr-nats"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
)

type EventRouter struct {
	app        *gofr.App
	natsClient *nats.PubSubWrapper
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
	}
}

func (er *EventRouter) Start() {
	er.app.Subscribe("raw_events", er.handleRawEvent)
	er.app.Run()
}

func (er *EventRouter) handleRawEvent(c *gofr.Context) error {
	var rawEvent map[string]interface{}
	if err := c.Bind(&rawEvent); err != nil {
		return err
	}

	event := cloudevents.NewEvent()
	event.SetID(uuid.New().String())
	event.SetSource("event-router")

	eventType, ok := rawEvent["type"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid event type")
	}
	event.SetType(eventType)

	tenantID, ok := rawEvent["tenant_id"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid tenant_id")
	}
	event.SetExtension("tenantid", tenantID)

	err := event.SetData(cloudevents.ApplicationJSON, rawEvent)
	if err != nil {
		return err
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Route to appropriate consumer queue based on event type
	consumerQueue := "events." + eventType
	if err := er.natsClient.Publish(c.Context, consumerQueue, eventJSON); err != nil {
		return err
	}

	return nil
}

func main() {
	router := NewEventRouter()
	router.Start()
}
