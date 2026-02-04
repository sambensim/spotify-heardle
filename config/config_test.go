// Package config handles application configuration from environment variables.
package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test_client_id")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test_secret")
	t.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback")
	t.Setenv("SESSION_SECRET", "test_session_secret")
	t.Setenv("PORT", "3000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.SpotifyClientID != "test_client_id" {
		t.Errorf("SpotifyClientID = %q, want %q", cfg.SpotifyClientID, "test_client_id")
	}

	if cfg.SpotifyClientSecret != "test_secret" {
		t.Errorf("SpotifyClientSecret = %q, want %q", cfg.SpotifyClientSecret, "test_secret")
	}

	if cfg.SpotifyRedirectURI != "http://localhost:8080/callback" {
		t.Errorf("SpotifyRedirectURI = %q, want %q", cfg.SpotifyRedirectURI, "http://localhost:8080/callback")
	}

	if cfg.SessionSecret != "test_session_secret" {
		t.Errorf("SessionSecret = %q, want %q", cfg.SessionSecret, "test_session_secret")
	}

	if cfg.Port != "3000" {
		t.Errorf("Port = %q, want %q", cfg.Port, "3000")
	}
}

func TestLoadDefaults(t *testing.T) {
	t.Setenv("SPOTIFY_CLIENT_ID", "test_id")
	t.Setenv("SPOTIFY_CLIENT_SECRET", "test_secret")
	t.Setenv("SPOTIFY_REDIRECT_URI", "http://localhost:8080/callback")
	t.Setenv("SESSION_SECRET", "test_session")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want default %q", cfg.Port, "8080")
	}
}

func TestLoadMissingRequired(t *testing.T) {
	tests := []struct {
		name       string
		envs       map[string]string
		wantErrMsg string
	}{
		{
			name: "missing client ID",
			envs: map[string]string{
				"SPOTIFY_CLIENT_SECRET": "secret",
				"SPOTIFY_REDIRECT_URI":  "http://localhost:8080/callback",
				"SESSION_SECRET":        "session",
			},
			wantErrMsg: "SPOTIFY_CLIENT_ID",
		},
		{
			name: "missing client secret",
			envs: map[string]string{
				"SPOTIFY_CLIENT_ID":    "id",
				"SPOTIFY_REDIRECT_URI": "http://localhost:8080/callback",
				"SESSION_SECRET":       "session",
			},
			wantErrMsg: "SPOTIFY_CLIENT_SECRET",
		},
		{
			name: "missing redirect URI",
			envs: map[string]string{
				"SPOTIFY_CLIENT_ID":     "id",
				"SPOTIFY_CLIENT_SECRET": "secret",
				"SESSION_SECRET":        "session",
			},
			wantErrMsg: "SPOTIFY_REDIRECT_URI",
		},
		{
			name: "missing session secret",
			envs: map[string]string{
				"SPOTIFY_CLIENT_ID":     "id",
				"SPOTIFY_CLIENT_SECRET": "secret",
				"SPOTIFY_REDIRECT_URI":  "http://localhost:8080/callback",
			},
			wantErrMsg: "SESSION_SECRET",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}

			_, err := Load()
			if err == nil {
				t.Fatal("Load() succeeded, want error")
			}
		})
	}
}
