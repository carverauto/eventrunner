package eventrunner

import (
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

func (a *AppWrapper) AddPubSub(pubsubClient container.PubSubProvider) {
	a.app.AddPubSub(pubsubClient)
}

func (a *AppWrapper) AddCassandra(cassandraClient container.CassandraProvider) {
	a.app.AddCassandra(cassandraClient)
}

func (a *AppWrapper) Migrate(migrationsMap map[int64]migration.Migrate) {
	a.app.Migrate(migrationsMap)
}
