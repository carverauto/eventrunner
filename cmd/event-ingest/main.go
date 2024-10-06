package main

import (
	"context"
	"errors"
	"log"

	"github.com/carverauto/eventrunner/pkg/api/middleware"
	"github.com/carverauto/eventrunner/pkg/config"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	"github.com/carverauto/eventrunner/pkg/eventingest"
	"gofr.dev/pkg/gofr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	app := gofr.New()

	// Load OAuth configuration
	oauthConfig := config.LoadOAuthConfig(app)

	// Initialize JWT middleware
	jwtMiddleware, err := middleware.NewJWTMiddleware(context.Background(), oauthConfig)
	if err != nil {
		log.Fatalf("Failed to initialize JWT middleware: %v", err)
	}

	app.UseMiddleware(middleware.CustomHeadersMiddleware())

	// Set up gRPC connection to API
	grpcServerAddress := app.Config.Get("GRPC_SERVER_ADDRESS")

	conn, err := grpc.Dial(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create gRPC event forwarder
	eventForwarder := eventingest.NewGRPCEventForwarder(conn)

	// Create and set up HTTP server
	httpServer := eventingest.NewHTTPServer(app, eventForwarder)

	// Register routes with middleware chain
	app.POST("/events", combineMiddleware(
		jwtMiddleware.Validate,
		middleware.AuthenticateAPIKey,
		middleware.RequireRole("admin", "event_publisher"),
		func(cc customctx.Context) (interface{}, error) {
			return httpServer.HandleEvent(cc)
		},
	))

	// Run the application
	app.Run()
}

// combineMiddleware chains multiple middleware functions together
func combineMiddleware(middlewares ...interface{}) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		// Retrieve the custom context from the original context
		customCtx, ok := c.Request.Context().Value("customCtx").(*customctx.CustomContext)
		if !ok {
			return nil, errors.New("failed to retrieve custom context")
		}

		// Define the final handler that will be called after applying all middleware
		finalHandler := func(ctx customctx.Context) (interface{}, error) {
			return nil, eventingest.NewInternalError("No handler provided")
		}

		// Apply middlewares in reverse order to build the middleware chain
		for i := len(middlewares) - 1; i >= 0; i-- {
			switch m := middlewares[i].(type) {
			case func(customctx.Context) (interface{}, error):
				// Set the final handler to the current one if no other handler is set
				if i == len(middlewares)-1 {
					finalHandler = m
				} else {
					// Wrap the final handler in the current function
					finalHandler = func(ctx customctx.Context) (interface{}, error) {
						return m(ctx)
					}
				}
			case func(func(customctx.Context) (interface{}, error)) func(customctx.Context) (interface{}, error):
				// Wrap the final handler in middleware if it's a middleware function
				finalHandler = m(finalHandler)
			}
		}

		// Execute the final middleware chain with the custom context
		return finalHandler(customCtx)
	}
}
