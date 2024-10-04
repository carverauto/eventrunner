// Package eventrunner pkg/eventrunner/interfaces.go
package eventrunner

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
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

type NATSClient interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(topic string, handler gofr.SubscribeFunc)
	Connect() error
}

type AppInterface interface {
	Subscribe(topic string, handler gofr.SubscribeFunc)
	Run()
	Logger() logging.Logger
	Metrics() metrics.Manager
	AddPubSub(pubsub container.PubSubProvider)
	AddCassandra(cassandraClient container.CassandraProvider)
	Migrate(migrationsMap map[int64]migration.Migrate)
}
