// Package config handles application configuration from environment variables.
package config

import (
	"fmt"
	"os"
)

// Config holds application configuration.
type Config struct {
	SpotifyClientID     string
	SpotifyClientSecret string
	SpotifyRedirectURI  string
	SessionSecret       string
	Port                string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_ID is required")
	}

	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_SECRET is required")
	}

	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")
	if redirectURI == "" {
		return nil, fmt.Errorf("SPOTIFY_REDIRECT_URI is required")
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		return nil, fmt.Errorf("SESSION_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		SpotifyClientID:     clientID,
		SpotifyClientSecret: clientSecret,
		SpotifyRedirectURI:  redirectURI,
		SessionSecret:       sessionSecret,
		Port:                port,
	}, nil
}
