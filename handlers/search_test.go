// Package handlers provides HTTP request handlers.
package handlers

import (
	"net/http"
	"net/http/httptest"
	"spotify-heardle/config"
	"spotify-heardle/storage"
	"testing"
)

func TestNewSearchHandler(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)

	handler := NewSearchHandler(authHandler)

	if handler == nil {
		t.Fatal("NewSearchHandler() returned nil")
	}
}

func TestHandleSearchNoAuth(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)
	handler := NewSearchHandler(authHandler)

	req := httptest.NewRequest("GET", "/api/search?q=test", nil)
	w := httptest.NewRecorder()

	handler.HandleSearch(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandleSearchMissingQuery(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)
	handler := NewSearchHandler(authHandler)

	req := httptest.NewRequest("GET", "/api/search", nil)
	w := httptest.NewRecorder()

	handler.HandleSearch(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d (unauthorized because no session)", w.Code, http.StatusUnauthorized)
	}
}
