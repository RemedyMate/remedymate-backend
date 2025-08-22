package config

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// OAuthConfig holds configurations for Google and Facebook OAuth providers
type OAuthConfig struct {
	Google   *oauth2.Config
	Facebook *oauth2.Config
}

// LoadOAuthConfig loads OAuth configurations from environment variables
func LoadOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		Google: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		Facebook: &oauth2.Config{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
			Scopes:       []string{"email", "public_profile"},
			Endpoint:     facebook.Endpoint,
		},
	}
}

// ValidateConfig checks if all required OAuth configurations are present
func (c *OAuthConfig) ValidateConfig() error {
	providers := map[string]struct {
		clientID     string
		clientSecret string
		redirectURL  string
	}{
		"Google":   {c.Google.ClientID, c.Google.ClientSecret, c.Google.RedirectURL},
		"Facebook": {c.Facebook.ClientID, c.Facebook.ClientSecret, c.Facebook.RedirectURL},
	}

	for name, config := range providers {
		if config.clientID == "" || config.clientSecret == "" || config.redirectURL == "" {
			return fmt.Errorf("%s OAuth configuration is incomplete", name)
		}
	}
	return nil
}
