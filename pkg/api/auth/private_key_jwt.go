package auth

import (
	"crypto/rsa"
	"time"

	"github.com/carverauto/eventrunner/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	ory "github.com/ory/client-go"
)

func createClientAssertion(clientID string, audience string, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss": clientID,
		"sub": clientID,
		"aud": audience,
		"jti": uuid.New().String(),
		"exp": now.Add(time.Minute * 10).Unix(),
		"iat": now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}

func InitializeOryClient(oryConfig *config.OAuthConfig, privateKey *rsa.PrivateKey) *ory.APIClient {
	configuration := ory.NewConfiguration()
	configuration.Servers = []ory.ServerConfiguration{
		{URL: oryConfig.OryProjectURL},
	}

	// Create client assertion
	clientAssertion, err := createClientAssertion(oryConfig.ClientID, oryConfig.OryProjectURL+"/oauth2/token", privateKey)
	if err != nil {
		panic(err)
	}

	// Set up the client to use JWT authentication
	configuration.AddDefaultHeader("Authorization", "Bearer "+clientAssertion)

	return ory.NewAPIClient(configuration)
}
