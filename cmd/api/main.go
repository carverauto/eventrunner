package main

import (
	"github.com/carverauto/eventrunner/cmd/api/migrations"
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/mongo"
)

func main() {
	app := gofr.New()

	// Set up MongoDB
	db := mongo.New(mongo.Config{URI: "mongodb://localhost:27017", Database: "eventrunner"})
	app.AddMongo(db)

	// Run migrations
	app.Migrate(migrations.All())

	// Set up routes
	tenantHandler := &handlers.TenantHandler{}
	userHandler := &handlers.UserHandler{}

	// Tenant routes (protected by API key)
	app.POST("/tenants", tenantHandler.Create, middleware.AuthenticateAPIKey)
	app.GET("/tenants", tenantHandler.GetAll, middleware.AuthenticateAPIKey)

	// User routes (protected by API key and role-based access)
	app.POST("/tenants/{tenant_id}/users", userHandler.Create, middleware.AuthenticateAPIKey, middleware.RequireRole("admin"))
	app.GET("/tenants/{tenant_id}/users", userHandler.GetAll, middleware.AuthenticateAPIKey, middleware.RequireRole("admin", "user"))

	// Run the application
	app.Run()
}
