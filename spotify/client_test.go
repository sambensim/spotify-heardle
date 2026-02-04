// Package spotify provides Spotify API client functionality.
package spotify

import (
	"spotify-heardle/models"
	"testing"
)

func TestNewClient(t *testing.T) {
	token := &models.Token{
		AccessToken:  "test_token",
		RefreshToken: "refresh_token",
	}

	client := NewClient(token)

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.token != token {
		t.Error("token not set correctly")
	}
}

func TestGetUserProfile(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}

func TestGetUserPlaylists(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}

func TestGetPlaylistTracks(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}

func TestSearchTracks(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}

func TestMakeRequest(t *testing.T) {
	t.Skip("Integration test - requires mock HTTP server")
}
