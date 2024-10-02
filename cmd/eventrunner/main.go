package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	nats "github.com/carverauto/gofr-nats"
	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	subjects := strings.Split(",", os.Getenv("NATS_SUBJECTS"))

	natsClient := nats.New(&nats.Config{
		Server: os.Getenv("PUBSUB_BROKER"),
		Stream: nats.StreamConfig{
			Stream: os.Getenv("NATS_STREAM"),
			// Subjects: []string{"orders.*", "shipments.*"},
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

	app.Subscribe("events.products", func(c *gofr.Context) error {
		c.Logger.Debug("Received message on 'products' subject")

		var productInfo struct {
			ProductID int     `json:"productId"`
			Price     float64 `json:"price"`
		}

		err := c.Bind(&productInfo)
		if err != nil {
			c.Logger.Error("Error binding product message:", err)
			return nil
		}

		c.Logger.Info("Received product", productInfo)

		return nil
	})

	app.Subscribe("events.order-logs", func(c *gofr.Context) error {
		c.Logger.Debug("Received message on 'order-logs' subject")

		var orderStatus struct {
			OrderID string `json:"orderId"`
			Status  string `json:"status"`
		}

		err := c.Bind(&orderStatus)
		if err != nil {
			c.Logger.Error("Error binding order message:", err)
			return nil
		}

		c.Logger.Info("Received order", orderStatus)

		return nil
	})

	fmt.Println("Subscribing to 'products' and 'order-logs' subjects...")
	app.Run()
}
