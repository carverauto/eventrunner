package main

import (
	"os"

	"gofr.dev/pkg/gofr"
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

	return string(data), nil
}

func RulesHandler(c *gofr.Context) (interface{}, error) {
	data, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, err
	}

	return string(data), nil
}
