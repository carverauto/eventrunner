package main

import (
	"os"

	"github.com/carverauto/eventrunner/cmd/api/migrations"
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
	ory "github.com/ory/client-go"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/mongo"
)

func main() {
	app := gofr.New()

	// Set up MongoDB
	db := mongo.New(mongo.Config{URI: os.Getenv("DB_URL"), Database: os.Getenv("DB_NAME")})
	app.AddMongo(db)

	// Run migrations
	app.Migrate(migrations.All())

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
