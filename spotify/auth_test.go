// Package spotify provides Spotify API client functionality.
package spotify

import (
	"spotify-heardle/models"
	"strings"
	"testing"
	"time"
)

func TestNewAuthManager(t *testing.T) {
	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURI := "http://localhost:8080/callback"

	auth := NewAuthManager(clientID, clientSecret, redirectURI)

	if auth == nil {
		t.Fatal("NewAuthManager() returned nil")
	}

	if auth.clientID != clientID {
		t.Errorf("clientID = %q, want %q", auth.clientID, clientID)
	}

	if auth.clientSecret != clientSecret {
		t.Errorf("clientSecret = %q, want %q", auth.clientSecret, clientSecret)
	}

	if auth.redirectURI != redirectURI {
		t.Errorf("redirectURI = %q, want %q", auth.redirectURI, redirectURI)
	}
}

func TestGetAuthURL(t *testing.T) {
	auth := NewAuthManager("client_id", "client_secret", "http://localhost:8080/callback")
	state := "random_state"

	url := auth.GetAuthURL(state)

	if url == "" {
		t.Error("GetAuthURL() returned empty string")
	}

	expectedPrefix := "https://accounts.spotify.com/authorize"
	if len(url) < len(expectedPrefix) || url[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("URL doesn't start with expected prefix. Got: %s", url)
	}
	
	if !strings.Contains(url, "streaming") {
		t.Error("URL missing 'streaming' scope")
	}
	
	if !strings.Contains(url, "user-modify-playback-state") {
		t.Error("URL missing 'user-modify-playback-state' scope")
	}
}

func TestExchangeCodeForToken(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}

func TestRefreshToken(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}

func TestCalculateExpiresAt(t *testing.T) {
	auth := &AuthManager{}
	expiresIn := 3600

	before := time.Now()
	expiresAt := auth.calculateExpiresAt(expiresIn)
	after := time.Now().Add(time.Duration(expiresIn) * time.Second)

	if expiresAt.Before(before) || expiresAt.After(after.Add(1*time.Second)) {
		t.Errorf("expiresAt out of expected range: %v", expiresAt)
	}
}

func TestParseTokenResponse(t *testing.T) {
	auth := &AuthManager{}
	
	tokenResp := tokenResponse{
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_456",
		ExpiresIn:    3600,
	}

	token := auth.parseTokenResponse(tokenResp)

	if token.AccessToken != tokenResp.AccessToken {
		t.Errorf("AccessToken = %q, want %q", token.AccessToken, tokenResp.AccessToken)
	}

	if token.RefreshToken != tokenResp.RefreshToken {
		t.Errorf("RefreshToken = %q, want %q", token.RefreshToken, tokenResp.RefreshToken)
	}

	if token.ExpiresAt.IsZero() {
		t.Error("ExpiresAt is zero")
	}
}

func TestParseTokenResponseNoRefresh(t *testing.T) {
	auth := &AuthManager{}
	existingToken := &models.Token{
		AccessToken:  "old_access",
		RefreshToken: "existing_refresh",
		ExpiresAt:    time.Now(),
	}

	tokenResp := tokenResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "",
		ExpiresIn:    3600,
	}

	token := auth.parseTokenResponseWithExisting(tokenResp, existingToken)

	if token.AccessToken != tokenResp.AccessToken {
		t.Errorf("AccessToken = %q, want %q", token.AccessToken, tokenResp.AccessToken)
	}

	if token.RefreshToken != existingToken.RefreshToken {
		t.Errorf("RefreshToken = %q, want existing %q", token.RefreshToken, existingToken.RefreshToken)
	}
}
