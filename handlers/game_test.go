// Package handlers provides HTTP request handlers.
package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"spotify-heardle/config"
	"spotify-heardle/models"
	"spotify-heardle/storage"
	"testing"
)

func TestNewGameHandler(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)

	handler := NewGameHandler(authHandler, store)

	if handler == nil {
		t.Fatal("NewGameHandler() returned nil")
	}
}

func TestHandleStartGameNoAuth(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)
	handler := NewGameHandler(authHandler, store)

	body := bytes.NewBufferString(`{"playlistIds":["playlist123"]}`)
	req := httptest.NewRequest("POST", "/api/game/start", body)
	w := httptest.NewRecorder()

	handler.HandleStartGame(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestSelectRandomTrack(t *testing.T) {
	tracks := []models.Track{
		{ID: "track1", Name: "Song 1", PreviewURL: "http://preview1"},
		{ID: "track2", Name: "Song 2", PreviewURL: "http://preview2"},
		{ID: "track3", Name: "Song 3", PreviewURL: "http://preview3"},
	}

	selected := selectRandomTrack(tracks)

	if selected.ID == "" {
		t.Error("selectRandomTrack() returned empty track")
	}

	found := false
	for _, track := range tracks {
		if track.ID == selected.ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("selectRandomTrack() returned track not in list")
	}
}

func TestFilterTracksWithPreview(t *testing.T) {
	tracks := []models.Track{
		{ID: "track1", Name: "Song 1", PreviewURL: "http://preview1"},
		{ID: "track2", Name: "Song 2", PreviewURL: ""},
		{ID: "track3", Name: "Song 3", PreviewURL: "http://preview3"},
		{ID: "track4", Name: "Song 4", PreviewURL: ""},
	}

	filtered := filterTracksWithPreview(tracks)

	if len(filtered) != 2 {
		t.Errorf("len(filtered) = %d, want 2", len(filtered))
	}

	for _, track := range filtered {
		if track.PreviewURL == "" {
			t.Errorf("filtered track %s has no preview URL", track.ID)
		}
	}
}

func TestHandleSubmitGuessNoSession(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)
	handler := NewGameHandler(authHandler, store)

	body := bytes.NewBufferString(`{"sessionId":"session123","trackId":"track1","trackName":"Song"}`)
	req := httptest.NewRequest("POST", "/api/game/guess", body)
	w := httptest.NewRecorder()

	handler.HandleSubmitGuess(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestHandleSkipNoSession(t *testing.T) {
	cfg := &config.Config{
		SpotifyClientID:     "test_id",
		SpotifyClientSecret: "test_secret",
		SpotifyRedirectURI:  "http://localhost:8080/callback",
		SessionSecret:       "test_session_secret",
	}
	store := storage.NewMemoryStore()
	authHandler := NewAuthHandler(cfg, store)
	handler := NewGameHandler(authHandler, store)

	body := bytes.NewBufferString(`{"sessionId":"session123"}`)
	req := httptest.NewRequest("POST", "/api/game/skip", body)
	w := httptest.NewRecorder()

	handler.HandleSkip(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
