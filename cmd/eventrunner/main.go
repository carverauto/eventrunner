// Package main cmd/eventrunner/main.go
package main

import (
	"github.com/carverauto/eventrunner/pkg/eventrunner"
)

func main() {

	router := eventrunner.NewEventRouter()

	// Add any middleware if needed
	// router.Use(someMiddleware)

	router.Start()
}
