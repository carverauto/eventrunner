package main

import (
	"encoding/json"
	"os"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

const (
	rulesPath = "/app/config/rules.json"
)

func main() {
	app := gofr.New()

	app.GET("/rules.json", RulesHandler)

	app.Run()
}

func RulesHandler(c *gofr.Context) (interface{}, error) {
	data, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, err
	}

	var jsonObj interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, err
	}

	return response.Raw{Data: jsonObj}, nil
}
