package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

func generateJWKS(publicKeyPath string) (string, error) {
	// Read the public key file
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key file: %v", err)
	}

	// Parse the PEM encoded public key
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Parse the public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// Assert that the public key is an RSA public key
	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("public key is not an RSA public key")
	}

	// Create a JWK from the RSA public key
	key, err := jwk.FromRaw(rsaPublicKey)
	if err != nil {
		return "", fmt.Errorf("failed to create JWK: %v", err)
	}

	// Set the key ID (kid) - this is optional but recommended
	if err := key.Set(jwk.KeyIDKey, "eventrunner-jwt"); err != nil {
		return "", fmt.Errorf("failed to set key ID: %v", err)
	}

	// Create a JWK Set with our key
	set := jwk.NewSet()
	err = set.AddKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to add key to set: %v", err)
	}

	// Marshal the JWK Set to JSON
	buf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JWKS: %v", err)
	}

	return string(buf), nil
}

func main() {
	jwks, err := generateJWKS("./public_key.pem")
	if err != nil {
		fmt.Printf("Error generating JWKS: %v\n", err)
		return
	}

	fmt.Println("Generated JWKS:")
	fmt.Println(jwks)

	// You can now use this JWKS in your Ory client configuration
}
