// File: main.go

package main

import (
	"context"
	"time"

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
	db := mongo.New(&mongo.Config{URI: "mongodb://er-mongodb.mongo.svc.cluster.local:27017", Database: "eventrunner"})

	// setup a context with a timeout
	dbCtx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	err := app.AddMongo(dbCtx, db)
	if err != nil {
		app.Logger().Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Initialize Ory client
	/*
		oryClient := ory.NewConfiguration()
		oryClient.Servers = ory.ServerConfigurations{{URL: os.Getenv("ORY_SDK_URL")}}
		// oryClient.DefaultHeader["Authorization"] = "Bearer " + os.Getenv("ORY_PAT")

		apiClient := ory.NewAPIClient(oryClient)

		// Initialize handlers
		h := handlers.NewHandlers(apiClient)
	*/

	// Add debug logging middleware

	// Set up routes
	/*
		app.POST("/api/users", middleware.Adapt(h.CreateUser, middleware.RequireRole("superuser", "tenant_admin")))
		app.POST("/api/tenants/:tenant_id/users", middleware.Adapt(h.CreateUser, middleware.RequireRole("superuser", "tenant_admin")))

		app.GET("/api/tenants/:tenant_id/users",
			debugMiddleware(middleware.Adapt(h.GetAllUsers, middleware.RequireRole("superuser", "tenant_admin"))))
	*/

	// Add the debug headers middleware
	app.UseMiddleware(middleware.DebugHeadersMiddleware())

	// Add other middleware and routes
	app.UseMiddleware(middleware.CustomHeadersMiddleware())

	// Routes...
	// app.POST("/api/v1/users", handlers.CreateUser)
	// Add a callback handler to your routes
	app.GET("/callback", func(ctx *gofr.Context) (interface{}, error) {
		code := ctx.Request.Param("code")
		state := ctx.Request.Param("state")

		ctx.Logger.Infof("Received callback with code: %s and state: %s", code, state)

		return map[string]string{
			"code":  code,
			"state": state,
		}, nil
	})

	app.GET("/api/test", testEndpoint)

	// Run the application
	app.Run()
}

func testEndpoint(c *gofr.Context) (interface{}, error) {
	return map[string]string{"message": "Hello, world!"}, nil
}
