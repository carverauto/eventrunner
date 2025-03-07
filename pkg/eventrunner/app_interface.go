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

package eventrunner

import (
	"context"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/logging"
	"gofr.dev/pkg/gofr/metrics"
	"gofr.dev/pkg/gofr/migration"
)

// AppWrapper wraps a *gofr.App and implements AppInterface.
type AppWrapper struct {
	app *gofr.App
}

func (a *AppWrapper) AddMongo(ctx context.Context, mongoClient container.MongoProvider) error {
	return a.app.AddMongo(ctx, mongoClient)
}

func NewAppWrapper(app *gofr.App) *AppWrapper {
	return &AppWrapper{app: app}
}

func (a *AppWrapper) Subscribe(topic string, handler gofr.SubscribeFunc) {
	a.app.Subscribe(topic, handler)
}

func (a *AppWrapper) Run() {
	a.app.Run()
}

func (a *AppWrapper) Logger() logging.Logger {
	return a.app.Logger()
}

func (a *AppWrapper) Metrics() metrics.Manager {
	return a.app.Metrics()
}

func (a *AppWrapper) AddPubSub(ctx context.Context, pubsubClient container.PubSubProvider) error {
	return a.app.AddPubSub(ctx, pubsubClient)
}

func (a *AppWrapper) AddCassandra(ctx context.Context, cassandraClient container.CassandraProvider) error {
	return a.app.AddCassandra(ctx, cassandraClient)
}

func (a *AppWrapper) Migrate(migrationsMap map[int64]migration.Migrate) {
	a.app.Migrate(migrationsMap)
}
