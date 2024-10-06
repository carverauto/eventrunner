package main

import (
	"github.com/carverauto/eventrunner/cmd/api/migrations"
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	middlewarePkg "github.com/carverauto/eventrunner/pkg/api/middleware"
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
	app.POST("/tenants", adapt(tenantHandler.Create, middleware.AuthenticateAPIKey))
	app.GET("/tenants", adapt(tenantHandler.GetAll, middleware.AuthenticateAPIKey))

	// User routes (protected by API key and role-based access)
	app.POST("/tenants/{tenant_id}/users", adapt(userHandler.Create,
		middleware.AuthenticateAPIKey,
		middleware.RequireRole("admin")))
	app.GET("/tenants/{tenant_id}/users", adapt(userHandler.GetAll,
		middleware.AuthenticateAPIKey,
		middleware.RequireRole("admin", "user")))

	// Run the application
	app.Run()
}

// adapt converts a handler func and middlewares into a gofr.Handler
func adapt(h interface{}, middlewares ...handlers.Middleware) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		var handler handlers.Handler
		switch h := h.(type) {
		case func(*gofr.Context) (interface{}, error):
			handler = handlers.HandlerFunc(h)
		default:
			if h, ok := h.(handlers.Handler); ok {
				handler = h
			} else {
				panic("unsupported handler type")
			}
		}

		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}

		return handler.Handle(c)
	}
}
