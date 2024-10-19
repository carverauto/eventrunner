package main

import (
	"context"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
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

	// Enable OAuth
	app.EnableOAuth(os.Getenv("JWKS_SERVER"), 20)

	// Set up routes
	tenantHandler := &handlers.TenantHandler{}
	userHandler := &handlers.UserHandler{}

	// Tenant routes (protected by OAuth)
	app.POST("/tenants", middleware.Adapt(tenantHandler.Create, middleware.RequireRole("admin")))
	app.GET("/tenants", middleware.Adapt(tenantHandler.GetAll, middleware.RequireRole("admin", "user")))
	app.POST("/tenants/{tenant_id}/users", middleware.Adapt(userHandler.Create, middleware.RequireRole("admin")))
	app.GET("/tenants/{tenant_id}/users", middleware.Adapt(userHandler.GetAll, middleware.RequireRole("admin", "user")))

	app.POST("/superuser", userHandler.CreateSuperUser)
	app.POST("/users", userHandler.CreateUser)

	// Run the application
	app.Run()
}
