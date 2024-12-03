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
	router := eventrunner.NewEventRouter(ctx, appWrapper, natsClient, nil)

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
