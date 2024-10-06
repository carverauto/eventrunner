package config

import (
	"gofr.dev/pkg/gofr"
)

type OAuthConfig struct {
	KeycloakURL        string
	Realm              string
	ClientID           string
	ClientSecret       string
	TokenIntrospectURL string
}

func LoadOAuthConfig(app *gofr.App) *OAuthConfig {
	return &OAuthConfig{
		KeycloakURL:        app.Config.Get("KEYCLOAK_URL"),
		Realm:              app.Config.Get("KEYCLOAK_REALM"),
		ClientID:           app.Config.Get("OAUTH_CLIENT_ID"),
		ClientSecret:       app.Config.Get("OAUTH_CLIENT_SECRET"),
		TokenIntrospectURL: app.Config.Get("TOKEN_INTROSPECT_URL"),
	}
}
