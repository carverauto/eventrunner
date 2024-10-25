// File: main.go

package main

import (
	"gofr.dev/pkg/gofr"
)

func main() {
	// Initialize gofr application
	app := gofr.New()

	// Register the JWKS route
	app.GET("/jwks.json", JWKSHandler)

	// Run the application
	app.Run()
}

type JWK struct {
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func JWKSHandler(c *gofr.Context) (interface{}, error) {
	// Create the JWKS data
	jwks := JWKS{
		Keys: []JWK{
			{
				E:   "AQAB",
				Kid: "eventrunner-jwt",
				Kty: "RSA",
				N:   "viVXLTzUz5zrrTRFe59lc5JfjonbmnBxgGVD2RHG-FQXdKp-5xnuH5C9ZLujcew8jYoeFw6o7ab7PMONzru5UcjxadKXaC1uTId_chCDVVVSD80IlYtzgchhMBTpqZJY5hd6GybODwJj0ulcfpXmw43dF5CRC9uLbLuSvkVsELgcioUJnaMTZjisY9R5ApeUOLSAZGOacdlVBBZfQb8pVjBqJQQmcyzooLZdXq-hNvutnI15sPQLcoBXXat_n8lfrI2Jr_mlG_rcvAdhZXUGeu1NeWdJuaHFoHcbV-PeSnr0mAGZxFEdM6nFywqmjtiU3EXhDmqfrB7hMiWdbAueRQ",
			},
		},
	}

	// Return the JWKS data; gofr will marshal it to JSON
	return jwks, nil
}
