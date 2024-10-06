package main

import (
	nats "github.com/carverauto/gofr-nats"
	"os"
	"strings"
	"time"

	"github.com/carverauto/eventrunner/pkg/eventrunner"
	"gofr.dev/pkg/gofr"
)

func main() {
	// Initialize the Gofr app
	app := gofr.New()

	subjects := strings.Split(",", os.Getenv("NATS_SUBJECTS"))

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

	// Wrap the Gofr app in an AppWrapper to implement the AppInterface
	appWrapper := eventrunner.NewAppWrapper(app)

	// Initialize the EventRouter
	router := eventrunner.NewEventRouter(appWrapper, natsClient, nil)

	// Add any middleware if needed
	// router.Use(someMiddleware)

	// Start the router
	router.Start()
}
