package main

import "github.com/carverauto/eventrunner/pkg/eventrunner"

func main() {
	router := eventrunner.NewEventRouter()
	router.Start()
}
