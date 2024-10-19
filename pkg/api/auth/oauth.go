package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/carverauto/eventrunner/pkg/api/middleware"
	customctx "github.com/carverauto/eventrunner/pkg/context"
	ory "github.com/ory/client-go"
	"gofr.dev/pkg/gofr"
)

func HandleOAuthCallback(oryClient *ory.APIClient) gofr.Handler {
	return func(c *gofr.Context) (interface{}, error) {
		cc, ok := c.Request.Context().Value(middleware.CustomContextKey).(*customctx.CustomContext)
		if !ok {
			return nil, fmt.Errorf("failed to retrieve custom context")
		}

		// Extract the authorization code from the query parameters
		code := c.Request.Param("code")
		if code == "" {
			return nil, fmt.Errorf("no code found in callback")
		}

		// Create an authenticated context for Ory API calls
		oryAuthedContext := context.WithValue(c.Context, ory.ContextAccessToken, os.Getenv("ORY_API_KEY"))

		// Exchange the authorization code for tokens
		tokenResponse, r, err := oryClient.OAuth2API.Oauth2TokenExchange(oryAuthedContext).
			GrantType("authorization_code").
			Code(code).
			RedirectUri(c.Request.HostName() + "/auth/callback"). // Make sure this matches your Ory configuration
			Execute()
		if err != nil {
			return nil, fmt.Errorf("error exchanging code for token: %v\nFull HTTP response: %v", err, r)
		}

		// Store the tokens securely (this is just an example, you should use a more secure method)
		setSecureCookie(c, "access_token", *tokenResponse.AccessToken, 3600) // 1 hour expiry for access token
		if tokenResponse.RefreshToken != nil {
			setSecureCookie(c, "refresh_token", *tokenResponse.RefreshToken, 86400*30) // 30 days expiry for refresh token
		}

		// Store session information
		sessionData := map[string]interface{}{
			"access_token":  tokenResponse.AccessToken,
			"refresh_token": tokenResponse.RefreshToken,
			"token_type":    tokenResponse.TokenType,
			"expiry":        time.Now().Add(time.Second * time.Duration(*tokenResponse.ExpiresIn)),
		}
		sessionJSON, _ := json.Marshal(sessionData)
		setSecureCookie(c, "session", string(sessionJSON), 3600) // 1 hour expiry for session

		// Redirect the user to your application's main page or dashboard
		c.Header().Set("Location", "/dashboard") // Adjust this to your application's needs
		c.SetStatusCode(http.StatusFound)
		return nil, nil
	}
}

func setSecureCookie(c *gofr.Context, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   maxAge,
	}
	http.SetCookie(c.Writer(), cookie)
}
