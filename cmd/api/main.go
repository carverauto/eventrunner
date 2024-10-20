// File: main.go

package main

import (
	"context"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
	ory "github.com/ory/client-go"
	"gofr.dev/pkg/gofr"
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

	// Enable OAuth
	app.EnableOAuth(os.Getenv("JWKS_SERVER"), 20)

	// Set up routes
	app.POST("/users", middleware.Adapt(h.CreateUser, middleware.RequireRole("superuser", "tenant_admin")))
	app.POST("/tenants/:tenant_id/users", middleware.Adapt(h.CreateUser, middleware.RequireRole("superuser", "tenant_admin")))
	app.GET("/tenants/{tenant_id}/users", middleware.Adapt(h.GetAllUsers, middleware.RequireRole("superuser", "tenant_admin")))

	// Run the application
	app.Run()
}
