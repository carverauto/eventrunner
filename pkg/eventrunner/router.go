// Package eventrunner pkg/eventrunner/eventrunner.go
package eventrunner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/carverauto/eventrunner/cmd/eventrunner/migrations"
	nats "github.com/carverauto/gofr-nats"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"gofr.dev/pkg/gofr"
	cassandraPkg "gofr.dev/pkg/gofr/datasource/cassandra"
)

type EventRouter struct {
	app             *gofr.App
	natsClient      *nats.PubSubWrapper
	bufferPool      *sync.Pool
	middlewares     []Middleware
	consumerManager *ConsumerManager
}

type Middleware func(HandlerFunc) HandlerFunc
type HandlerFunc func(*gofr.Context, *cloudevents.Event) error

func NewEventRouter() *EventRouter {
	app := gofr.New()

	// Configure Cassandra
	cassandraConfig := cassandraPkg.Config{
		Hosts:    os.Getenv("CASSANDRA_HOSTS"),
		Keyspace: os.Getenv("CASSANDRA_KEYSPACE"),
		Port:     9042, // Default Cassandra port, adjust if necessary
		Username: os.Getenv("CASSANDRA_USERNAME"),
		Password: os.Getenv("CASSANDRA_PASSWORD"),
	}
	cassandra := cassandraPkg.New(cassandraConfig)
	app.AddCassandra(cassandra)

	// Add migrations to run
	app.Migrate(migrations.All())

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

	consumerManager := NewConsumerManager(app)
	cassandraSink := NewCassandraEventSink()

	consumerManager.AddConsumer("clickhouse", cassandraSink)

	return &EventRouter{
		app:             app,
		natsClient:      natsClient,
		bufferPool:      &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},
		consumerManager: consumerManager,
	}
}

func (er *EventRouter) Use(middleware Middleware) {
	er.middlewares = append(er.middlewares, middleware)
}

func (er *EventRouter) Start() {
	er.app.Subscribe("events.products", er.handleEvent)
	er.app.Run()
}

func (er *EventRouter) handleEvent(c *gofr.Context) error {
	// First, try to unmarshal as a CloudEvent
	var event cloudevents.Event
	err := c.Bind(&event)

	if err != nil {
		// If it's not a CloudEvent, treat it as a raw JSON message
		var rawMessage map[string]interface{}
		if err := c.Bind(&rawMessage); err != nil {
			return fmt.Errorf("failed to parse message: %w", err)
		}

		// Convert the raw message to a CloudEvent
		event = cloudevents.NewEvent()
		event.SetID(uuid.New().String())
		event.SetSource("eventrunner")
		event.SetType("com.example.event") // You might want to determine this based on the topic
		event.SetTime(time.Now())
		if err := event.SetData(cloudevents.ApplicationJSON, rawMessage); err != nil {
			return fmt.Errorf("failed to set event data: %w", err)
		}
	}

	// Apply middlewares
	handler := er.applyMiddleware(er.routeEvent)
	return handler(c, &event)
}

func (er *EventRouter) applyMiddleware(handler HandlerFunc) HandlerFunc {
	for i := len(er.middlewares) - 1; i >= 0; i-- {
		handler = er.middlewares[i](handler)
	}
	return handler
}

func (er *EventRouter) routeEvent(c *gofr.Context, event *cloudevents.Event) error {
	// Log event using ConsumerManager
	if err := er.consumerManager.ConsumeEvent(c, event); err != nil {
		er.app.Logger().Errorf("Failed to consume event: %v", err)
		// Continue processing even if logging fails
	}

	eventType := event.Type()
	tenantID, _ := event.Context.GetExtension("tenantid")

	// Route to appropriate consumer queue based on event type
	consumerQueue := fmt.Sprintf("events.%s", eventType)

	buf := er.getBuffer()
	defer er.putBuffer(buf)

	if err := json.NewEncoder(buf).Encode(event); err != nil {
		return fmt.Errorf("failed to encode event: %w", err)
	}

	messageID := uuid.New().String()
	if err := er.natsClient.Publish(c, consumerQueue, buf.Bytes()); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	er.app.Logger().Logf("published event %s to %s for tenant %s", messageID, consumerQueue, tenantID)

	return nil
}
