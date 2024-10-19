package main

import (
	"context"
	"crypto/rsa"
	"log"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/auth"
	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
	"github.com/carverauto/eventrunner/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/mongo"
)

const (
	dbConnectTimeout = 10 * time.Second
)

var privateKey *rsa.PrivateKey

func init() {
	// Load your private key
	privateKeyPEM, err := os.ReadFile("./private_key.pem")
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
}

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

	// Run migrations
	// app.Migrate(migrations.All())

	// Set up Ory client
	oryConfig := config.LoadOAuthConfig(app)
	oryClient := auth.InitializeOryClient(oryConfig, privateKey)

	// Set up routes
	tenantHandler := &handlers.TenantHandler{}
	userHandler := &handlers.UserHandler{}

	// Tenant routes (protected by Ory Auth)
	app.POST("/tenants", middleware.Adapt(tenantHandler.Create, middleware.OryAuthMiddleware(oryClient)))
	app.GET("/tenants", middleware.Adapt(tenantHandler.GetAll, middleware.OryAuthMiddleware(oryClient)))
	app.POST("/tenants/{tenant_id}/users", middleware.Adapt(userHandler.Create,
		middleware.OryAuthMiddleware(oryClient),
		middleware.RequireRole("admin")))
	app.GET("/tenants/{tenant_id}/users", middleware.Adapt(userHandler.GetAll,
		middleware.OryAuthMiddleware(oryClient),
		middleware.RequireRole("admin", "user")))

	app.GET("/auth/callback", auth.HandleOAuthCallback(oryClient))

	// Run the application
	app.Run()
}
