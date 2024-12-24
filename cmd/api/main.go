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
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/handlers"
	"github.com/carverauto/eventrunner/pkg/api/middleware"
	ory "github.com/ory/client-go"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
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
	// db := mongo.New(mongo.Config{URI: "mongodb://er-mongodb.mongo.svc.cluster.local:27017", Database: "eventrunner"})
	db := mongo.New(mongo.Config{URI: os.Getenv("MONGO_DSN"), Database: "eventrunner"})

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

	// lets test getting the TLS config for MongoDB
	tlsConfig, err := getTLSConfigForMongo()
	if err != nil {
		app.Logger().Errorf("Failed to get TLS config for MongoDB: %v", err)
		return
	}

	// print the TLS config
	app.Logger().Infof("TLS Config: %v", tlsConfig)

	// Run the application
	app.Run()
}

func testEndpoint(c *gofr.Context) (interface{}, error) {
	return map[string]string{"message": "Hello, world!"}, nil
}

func getTLSConfigForMongo() (*tls.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1) Create an X.509 Source, specifying the socket path via WithClientOptions
	source, err := workloadapi.NewX509Source(ctx,
		workloadapi.WithClientOptions(
			workloadapi.WithAddr("unix:///tmp/spire-agent.sock"), // your spire agent path
		),
	)
	if err != nil {
		return nil, err
	}
	// Optional: You may want to defer source.Close() if the source is long-lived
	defer source.Close()

	// 2) Fetch the SPIFFE X.509 SVID (certificate chain + private key)
	svid, err := source.GetX509SVID()
	if err != nil {
		return nil, err
	}

	// 3) Convert svid.Certificates ( []*x509.Certificate ) to raw [][]byte
	//    This is what tls.Certificate.Certificate expects.
	var certChain [][]byte
	for _, c := range svid.Certificates {
		certChain = append(certChain, c.Raw)
	}

	// 4) Build a tls.Certificate from the SVID
	x509Cert := tls.Certificate{
		Certificate: certChain,
		PrivateKey:  svid.PrivateKey,
	}

	// 5) Create a tls.Config that will present the SPIFFE certificate to MongoDB
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{x509Cert},
		InsecureSkipVerify: false, // typically 'false' in production
		MinVersion:         tls.VersionTLS13,
		// RootCAs:          <— Set if you also want to verify the Mongo server's cert
	}

	return tlsConfig, nil
}
