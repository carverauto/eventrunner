package config

import (
	"gofr.dev/pkg/gofr"
)

type OAuthConfig struct {
	OryProjectURL string
	ClientID      string
	RedirectURI   string
}

func LoadOAuthConfig(app *gofr.App) *OAuthConfig {
	return &OAuthConfig{
		OryProjectURL: app.Config.Get("ORY_PROJECT_URL"),
		ClientID:      app.Config.Get("ORY_CLIENT_ID"),
		RedirectURI:   app.Config.Get("ORY_REDIRECT_URI"),
	}
}
