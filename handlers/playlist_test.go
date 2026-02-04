// Package handlers provides HTTP request handlers.
package handlers

import (
	"net/http"
	"net/http/httptest"
	"spotify-heardle/config"
	"spotify-heardle/storage"
	"testing"
)

func TestNewPlaylistHandler(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)

	handler := NewPlaylistHandler(authHandler)

	if handler == nil {
		t.Fatal("NewPlaylistHandler() returned nil")
	}
}

func TestHandleGetPlaylistsNoAuth(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)
	handler := NewPlaylistHandler(authHandler)

	req := httptest.NewRequest("GET", "/api/playlists", nil)
	w := httptest.NewRecorder()

	handler.HandleGetPlaylists(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
