// Package main cmd/api-admin/main.go

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
	db := mongo.New(&mongo.Config{URI: os.Getenv("DB_URL"), Database: "eventrunner"})

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

	apiClient := ory.NewAPIClient(oryClient)

	// Initialize handlers
	h := handlers.NewHandlers(apiClient)

	// Add debug logging middleware

	app.UseMiddleware(middleware.DebugHeadersMiddleware())

	// Add other middleware and routes
	app.UseMiddleware(middleware.CustomHeadersMiddleware())

	// API Credentials routes
	app.POST("/api/credentials", middleware.Adapt(
		h.CreateAPICredential,
		middleware.RequireUser,
	))

	app.GET("/api/credentials", middleware.Adapt(
		h.ListAPICredentials,
		middleware.RequireUser,
	))

	// this endpoint is used by the Ory Kratos login flow
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
