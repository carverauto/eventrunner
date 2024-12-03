/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

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
	Connect(ctx context.Context) error
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
	AddPubSub(ctx context.Context, pubsub container.PubSubProvider) error
	AddCassandra(ctx context.Context, cassandraClient container.CassandraProvider) error
	AddMongo(ctx context.Context, mongoClient container.MongoProvider) error
	Migrate(migrationsMap map[int64]migration.Migrate)
}
