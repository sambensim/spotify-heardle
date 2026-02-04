// Package spotify provides Spotify API client functionality.
package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"spotify-heardle/models"
	"strings"
	"time"
)

const (
	authURL  = "https://accounts.spotify.com/authorize"
	tokenURL = "https://accounts.spotify.com/api/token"
)

// AuthManager handles Spotify OAuth authentication.
type AuthManager struct {
	clientID     string
	clientSecret string
	redirectURI  string
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// NewAuthManager creates a new Spotify auth manager.
func NewAuthManager(clientID, clientSecret, redirectURI string) *AuthManager {
	return &AuthManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

// GetAuthURL generates the Spotify authorization URL.
func (a *AuthManager) GetAuthURL(state string) string {
	scopes := []string{
		"user-read-private",
		"playlist-read-private",
		"playlist-read-collaborative",
		"streaming",
		"user-read-email",
		"user-modify-playback-state",
		"user-read-playback-state",
	}

	params := url.Values{}
	params.Set("client_id", a.clientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", a.redirectURI)
	params.Set("scope", strings.Join(scopes, " "))
	params.Set("state", state)

	return authURL + "?" + params.Encode()
}

// ExchangeCodeForToken exchanges authorization code for access token.
func (a *AuthManager) ExchangeCodeForToken(code string) (*models.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", a.redirectURI)
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	return a.parseTokenResponse(tokenResp), nil
}

// RefreshToken refreshes an expired access token.
func (a *AuthManager) RefreshToken(refreshToken string, existingToken *models.Token) (*models.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", a.clientID)
	data.Set("client_secret", a.clientSecret)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("refresh request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	return a.parseTokenResponseWithExisting(tokenResp, existingToken), nil
}

func (a *AuthManager) calculateExpiresAt(expiresIn int) time.Time {
	return time.Now().Add(time.Duration(expiresIn) * time.Second)
}

func (a *AuthManager) parseTokenResponse(resp tokenResponse) *models.Token {
	return &models.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    a.calculateExpiresAt(resp.ExpiresIn),
	}
}

func (a *AuthManager) parseTokenResponseWithExisting(resp tokenResponse, existing *models.Token) *models.Token {
	refreshToken := resp.RefreshToken
	if refreshToken == "" && existing != nil {
		refreshToken = existing.RefreshToken
	}

	return &models.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    a.calculateExpiresAt(resp.ExpiresIn),
	}
}
