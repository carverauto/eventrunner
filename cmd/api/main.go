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
	ctx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	err := app.AddMongo(ctx, db)
	if err != nil {
		app.Logger().Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Run migrations
	// app.Migrate(migrations.All())

	// Set up Ory client
	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{URL: os.Getenv("ORY_PROJECT_URL")},
	}
	oryClient := ory.NewAPIClient(configuration)

	// Set up routes
	tenantHandler := &handlers.TenantHandler{}
	userHandler := &handlers.UserHandler{}

	// Tenant routes (protected by Ory Auth)
	app.POST("/tenants", middleware.Adapt(tenantHandler.Create, middleware.OryAuthMiddleware(oryClient)))
	app.GET("/tenants", middleware.Adapt(tenantHandler.GetAll, middleware.OryAuthMiddleware(oryClient)))

	// User routes (protected by API key and role-based access)
	// TODO: Add Ory Auth middleware to user routes
	app.POST("/tenants/{tenant_id}/users", middleware.Adapt(userHandler.Create,
		middleware.AuthenticateAPIKey,
		middleware.RequireRole("admin")))
	app.GET("/tenants/{tenant_id}/users", middleware.Adapt(userHandler.GetAll,
		middleware.AuthenticateAPIKey,
		middleware.RequireRole("admin", "user")))

	// Run the application
	app.Run()
}
