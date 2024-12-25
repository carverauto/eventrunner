// File: main.go

/*
* Copyright 2024 Carver Automation Corp.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
*  limitations under the License.
 */

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
	ory "github.com/ory/client-go"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource/mongo"
)

const (
	dbConnectTimeout = 10 * time.Second
)

func main() {
	app := gofr.New()

	ctx := context.Background()

	app.Logger().Info("Starting EventRunner API")

	tlsConfig, err := getTLSConfigForMongo()
	if err != nil {
		app.Logger().Errorf("Failed to get TLS config for MongoDB: %v", err)
		return
	}

	config := mongo.Config{
		Database:          "eventrunner",
		ConnectionTimeout: dbConnectTimeout,
	}

	// Set up MongoDB
	clientOpts := options.Client().
		// ApplyURI("mongodb://er-mongodb.mongo.svc.cluster.local:27017").
		ApplyURI(os.Getenv("MONGO_DSN")).
		SetAuth(options.Credential{
			AuthMechanism: "MONGODB-X509",
			AuthSource:    "$external",
		}).
		SetTLSConfig(tlsConfig)

	client := mongo.New(config).WithClientOptions(clientOpts)
	client.UseLogger(app.Logger())
	client.UseMetrics(app.Metrics())

	err = client.Connect(ctx)
	if err != nil {
		app.Logger().Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// setup a context with a timeout
	dbCtx, cancel := context.WithTimeout(ctx, dbConnectTimeout)
	defer cancel()

	err = app.AddMongo(dbCtx, client)
	if err != nil {
		app.Logger().Errorf("Failed to connect to MongoDB: %v", err)
		return
	}

	// Initialize Ory client
	oryClient := ory.NewConfiguration()

	oryClient.Servers = ory.ServerConfigurations{{URL: "http://hydra-admin.auth:4445"}}

	apiClient := ory.NewAPIClient(oryClient)

	// Initialize handlers
	h := handlers.NewHandlers(apiClient)

	// Add debug logging middleware

	// Set up routes
	app.POST("/api/users", middleware.Adapt(h.CreateUser, middleware.RequireRole("superuser", "tenant_admin")))
	app.POST("/api/tenants/{tenant_id}/users", middleware.Adapt(h.CreateUser, middleware.RequireRole("superuser", "tenant_admin")))

	app.GET("/api/tenants/{tenant_id}/users",
		middleware.Adapt(h.GetAllUsers, middleware.RequireRole("superuser"), middleware.RequireTenant))

	// Add the debug headers middleware
	app.UseMiddleware(middleware.DebugHeadersMiddleware())

	// Add other middleware and routes
	app.UseMiddleware(middleware.CustomHeadersMiddleware())

	// API Credentials routes
	app.POST("/api/admin/credentials", middleware.Adapt(
		h.CreateAPICredential,
		middleware.RequireUser,
	))

	app.GET("/api/admin/credentials", middleware.Adapt(
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

func getTLSConfigForMongo() (*tls.Config, error) {
	log.Printf("Starting getTLSConfigForMongo")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Print the current process UID
	log.Printf("Current process UID: %d", os.Getuid())

	log.Printf("Creating X509Source with socket path: unix:///tmp/spire-agent.sock")
	source, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(
			workloadapi.WithAddr("unix:///tmp/spire-agent.sock"),
		),
	)
	if err != nil {
		log.Printf("Error creating X509Source: %v", err)
		return nil, fmt.Errorf("error creating X509Source: %v", err)
	}
	defer source.Close()

	log.Printf("Getting X509SVID")
	svid, err := source.GetX509SVID()
	if err != nil {
		log.Printf("Error getting X509SVID: %v", err)
		return nil, fmt.Errorf("error getting X509SVID: %v", err)
	}
	log.Printf("Got SVID with ID: %s", svid.ID)

	rootCAs := x509.NewCertPool()
	bundle, err := source.GetX509BundleForTrustDomain(spiffeid.RequireTrustDomainFromString("tunnel.threadr.ai"))
	if err != nil {
		return nil, err
	}
	for _, cert := range bundle.X509Authorities() {
		rootCAs.AddCert(cert)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{svid.Certificates[0].Raw},
				PrivateKey:  svid.PrivateKey,
			},
		},
		RootCAs:    rootCAs,
		MinVersion: tls.VersionTLS13,
	}, nil
}
