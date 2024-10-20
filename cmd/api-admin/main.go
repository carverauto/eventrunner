// File: main.go

package main

import (
	"context"
	"os"
	"slices"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	ory "github.com/ory/client-go"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/datasource/mongo"
)

const (
	dbConnectTimeout = 10 * time.Second
)

func main() {
	app := gofr.New()

	ctx := context.Background()

	// Set up MongoDB
	db := mongo.New(&mongo.Config{URI: "mongodb://localhost:27017", Database: "eventrunner"})

	// setup a context with a timeout
	dbCtx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	err := app.AddMongo(dbCtx, db)
	if err != nil {
		app.Logger().Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Initialize Ory client
	oryClient := ory.NewConfiguration()
	oryClient.Servers = ory.ServerConfigurations{{URL: os.Getenv("ORY_SDK_URL")}}
	oryClient.DefaultHeader["Authorization"] = "Bearer " + os.Getenv("ORY_PAT")

	apiClient := ory.NewAPIClient(oryClient)

	// Initialize handlers
	h := handlers.NewHandlers(apiClient)

	app.EnableAPIKeyAuthWithValidator(apiKeyValidator)

	// Set up routes
	app.POST("/superuser", h.CreateSuperUser)

	// Run the application
	app.Run()
}

func apiKeyValidator(_ *container.Container, apiKey string) bool {
	validKeys := []string{os.Getenv("ADMIN_API_KEY")}

	return slices.Contains(validKeys, apiKey)
}
