package config

import (
	"errors"
	"fmt"

	"github.com/hibare/Waypoint/internal/constants"
)

var (
	// ErrOIDCIssuerEmpty indicates that the OIDC issuer URL is empty.
	ErrOIDCIssuerEmpty = errors.New("oidc issuer url is empty")

	// ErrOIDCClientIDEmpty indicates that the OIDC client ID is empty.
	ErrOIDCClientIDEmpty = errors.New("oidc client id is empty")

	// ErrOIDCClientSecretEmpty indicates that the OIDC client secret is empty.
	ErrOIDCClientSecretEmpty = errors.New("oidc client secret is empty")
)

const (
	// OIDCScopeOpenID is the OIDC scope for OpenID Connect authentication.
	OIDCScopeOpenID = "openid"

	// OIDCScopeProfile is the OIDC scope for requesting access to the user's profile information.
	OIDCScopeProfile = "profile"

	// OIDCScopeEmail is the OIDC scope for requesting access to the user's email address.
	OIDCScopeEmail = "email"
)

// OIDCConfig holds OIDC authentication configuration.
type OIDCConfig struct {
	IssuerURL    string   `mapstructure:"issuer_url"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"` //nolint:gosec // required for OIDC authentication
	Scopes       []string `mapstructure:"scopes"`
}

// Validate checks if the OIDC configuration is valid.
func (o *OIDCConfig) Validate() error {
	if o.IssuerURL == "" {
		return ErrOIDCIssuerEmpty
	}
	if o.ClientID == "" {
		return ErrOIDCClientIDEmpty
	}
	if o.ClientSecret == "" {
		return ErrOIDCClientSecretEmpty
	}
	if len(o.Scopes) == 0 {
		o.Scopes = []string{OIDCScopeOpenID, OIDCScopeProfile, OIDCScopeEmail}
	}
	return nil
}

// GetRedirectURI returns the OIDC redirect URI constructed from BaseURL.
func (o *OIDCConfig) GetRedirectURI(baseURL string) string {
	return fmt.Sprintf("%s%s", baseURL, constants.OIDCCallbackPath)
}
