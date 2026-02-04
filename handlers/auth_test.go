// Package handlers provides HTTP request handlers.
package handlers

import (
	"net/http"
	"net/http/httptest"
	"spotify-heardle/config"
	"spotify-heardle/storage"
	"strings"
	"testing"
)

func TestNewAuthHandler(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()

	handler := NewAuthHandler(cfg, store)

	if handler == nil {
		t.Fatal("NewAuthHandler() returned nil")
	}
}

func TestHandleLogin(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	handler := NewAuthHandler(cfg, store)

	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	handler.HandleLogin(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("status = %d, want %d", w.Code, http.StatusTemporaryRedirect)
	}

	location := w.Header().Get("Location")
	if !strings.HasPrefix(location, "https://accounts.spotify.com/authorize") {
		t.Errorf("redirect location doesn't start with Spotify auth URL: %s", location)
	}
}

func TestGenerateState(t *testing.T) {
	state1, err := generateState()
	if err != nil {
		t.Fatalf("generateState() failed: %v", err)
	}

	if len(state1) == 0 {
		t.Error("generateState() returned empty string")
	}

	state2, err := generateState()
	if err != nil {
		t.Fatalf("generateState() failed: %v", err)
	}

	if state1 == state2 {
		t.Error("generateState() returned same state twice")
	}
}

func TestHandleGetTokenNoAuth(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	handler := NewAuthHandler(cfg, store)

	req := httptest.NewRequest("GET", "/api/token", nil)
	w := httptest.NewRecorder()

	handler.HandleGetToken(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestGenerateSessionID(t *testing.T) {
	id1, err := generateSessionID()
	if err != nil {
		t.Fatalf("generateSessionID() failed: %v", err)
	}

	if len(id1) == 0 {
		t.Error("generateSessionID() returned empty string")
	}

	id2, err := generateSessionID()
	if err != nil {
		t.Fatalf("generateSessionID() failed: %v", err)
	}

	if id1 == id2 {
		t.Error("generateSessionID() returned same ID twice")
	}
}
