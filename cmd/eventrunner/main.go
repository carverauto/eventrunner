package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/carverauto/eventrunner/pkg/eventrunner"
	nats "github.com/carverauto/gofr-nats"
	"gofr.dev/pkg/gofr"
)

func main() {
	// Initialize the Gofr app
	app := gofr.New()

	ctx := context.Background()

	subjects := strings.Split(app.Config.Get("NATS_SUBJECTS"), ",")

	natsClient := nats.New(&nats.Config{
		Server: os.Getenv("PUBSUB_BROKER"),
		Stream: nats.StreamConfig{
			Stream:   os.Getenv("NATS_STREAM"),
			Subjects: subjects,
			MaxBytes: 1024 * 1024 * 1024, // 1 GB, adjust as needed
		},
		MaxWait:     5 * time.Second,
		BatchSize:   100,
		MaxPullWait: 10,
		Consumer:    os.Getenv("NATS_CONSUMER"),
		CredsFile:   os.Getenv("NATS_CREDS_FILE"),
	}, app.Logger())
	natsClient.UseLogger(app.Logger)
	natsClient.UseMetrics(app.Metrics())
	if err := app.AddPubSub(ctx, natsClient); err != nil {
		app.Logger().Errorf("Failed to connect to NATS: %v", err)
		return
	}

	// Wrap the Gofr app in an AppWrapper to implement the AppInterface
	appWrapper := eventrunner.NewAppWrapper(app)

	// Initialize the EventRouter with the existing NATS client
	router := eventrunner.NewEventRouter(appWrapper, natsClient, nil)

	// Start the router in a goroutine
	go func() {
		router.Start()
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	app.Logger().Info("Shutting down gracefully...")
	// Perform any necessary cleanup here
	app.Logger().Info("Shutdown complete")
}
