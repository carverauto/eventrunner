package middleware

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
)

// OIDCVerifier is a wrapper around oidc.IDTokenVerifier
type OIDCVerifier struct {
	verifier *oidc.IDTokenVerifier
}

// NewOIDCVerifier returns an instance of OIDCVerifier
func NewOIDCVerifier(verifier *oidc.IDTokenVerifier) *OIDCVerifier {
	return &OIDCVerifier{verifier: verifier}
}

// Verify method to conform to the IDTokenVerifier interface
func (o *OIDCVerifier) Verify(ctx context.Context, rawToken string) (Token, error) {
	idToken, err := o.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, err
	}
	return &OIDCToken{idToken: idToken}, nil
}

// OIDCToken is a wrapper around oidc.IDToken to implement the Token interface
type OIDCToken struct {
	idToken *oidc.IDToken
}

// Claims wraps the oidc.IDToken.Claims method to conform to the Token interface
func (o *OIDCToken) Claims(v interface{}) error {
	return o.idToken.Claims(v)
}
