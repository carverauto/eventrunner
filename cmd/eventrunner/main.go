// Package main cmd/eventrunner/main.go
package main

import (
	"log"
	"os"

	"github.com/carverauto/eventrunner/pkg/eventrunner"
)

func main() {
	// Ensure environment variables are set
	requiredEnvVars := []string{
		"CASSANDRA_HOSTS",
		"CASSANDRA_KEYSPACE",
		"CASSANDRA_USERNAME",
		"CASSANDRA_PASSWORD",
		"NATS_SUBJECTS",
		"PUBSUB_BROKER",
		"NATS_STREAM",
		"NATS_CONSUMER",
		"NATS_CREDS_FILE",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Environment variable %s is not set", envVar)
		}
	}

	router := eventrunner.NewEventRouter()

	// Add any middleware if needed
	// router.Use(someMiddleware)

	router.Start()
}
