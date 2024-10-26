package main

import (
	"encoding/json"
	"os"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http/response"
)

const (
	jwksPath  = "/app/config/jwks.json"
	rulesPath = "/app/config/rules.json"
)

func main() {
	app := gofr.New()

	app.GET("/jwks.json", JWKSHandler)
	app.GET("/rules.json", RulesHandler)

	app.Run()
}

func JWKSHandler(c *gofr.Context) (interface{}, error) {
	data, err := os.ReadFile(jwksPath)
	if err != nil {
		return nil, err
	}

	var jsonObj interface{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return nil, err
	}

	return response.Raw{Data: jsonObj}, nil
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
