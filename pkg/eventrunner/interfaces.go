// Package eventrunner pkg/eventrunner/interfaces.go
package eventrunner

import (
	"context"
	"io"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/datasource"
	"gofr.dev/pkg/gofr/datasource/pubsub"
	"gofr.dev/pkg/gofr/logging"
	"gofr.dev/pkg/gofr/metrics"
	"gofr.dev/pkg/gofr/migration"
)

//go:generate mockgen -destination=mock_eventrunner.go -package=eventrunner -source=./interfaces.go EventSink,EventConsumer,NATSClient,AppInterface

type EventSink interface {
	LogEvent(context.Context, *cloudevents.Event) error
}

type EventConsumer interface {
	ConsumeEvent(ctx *gofr.Context, event *cloudevents.Event) error
}

type Buffer interface {
	io.Writer
	Bytes() []byte
	Reset()
}

type NATSClient interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string) (*pubsub.Message, error)
	Connect()
	Health() datasource.Health
	CreateTopic(ctx context.Context, name string) error
	DeleteTopic(ctx context.Context, name string) error
	Close() error
	UseLogger(logger any)
	UseMetrics(metrics any)
	UseTracer(tracer any)
}

type AppInterface interface {
	Subscribe(topic string, handler gofr.SubscribeFunc)
	Run()
	Logger() logging.Logger
	Metrics() metrics.Manager
	AddPubSub(pubsub container.PubSubProvider)
	AddCassandra(cassandraClient container.CassandraProvider)
	AddMongo(mongoClient container.MongoProvider)
	Migrate(migrationsMap map[int64]migration.Migrate)
}
