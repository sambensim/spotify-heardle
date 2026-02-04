// Package handlers provides HTTP request handlers.
package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"spotify-heardle/config"
	"spotify-heardle/models"
	"spotify-heardle/spotify"
	"spotify-heardle/storage"
)

// AuthHandler handles authentication routes.
type AuthHandler struct {
	auth  *spotify.AuthManager
	store *storage.MemoryStore
}

type sessionCookie struct {
	UserID string
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(cfg *config.Config, store *storage.MemoryStore) *AuthHandler {
	return &AuthHandler{
		auth:  spotify.NewAuthManager(cfg.SpotifyClientID, cfg.SpotifyClientSecret, cfg.SpotifyRedirectURI),
		store: store,
	}
}

// HandleLogin redirects to Spotify authorization page.
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := generateState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	authURL := h.auth.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleCallback handles the OAuth callback from Spotify.
func (h *AuthHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	token, err := h.auth.ExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	client := spotify.NewClient(token)
	profile, err := client.GetUserProfile()
	if err != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	user := models.NewUser(profile.ID, profile.DisplayName, token)
	if err := h.store.SaveUser(user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	sessionData := sessionCookie{UserID: user.ID}
	sessionJSON, _ := json.Marshal(sessionData)
	sessionValue := base64.StdEncoding.EncodeToString(sessionJSON)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    sessionValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   86400 * 7,
	})

	http.Redirect(w, r, "/playlists.html", http.StatusTemporaryRedirect)
}

// HandleLogout clears the session.
func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}

// GetUserFromSession retrieves user from session cookie.
func (h *AuthHandler) GetUserFromSession(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, fmt.Errorf("no session cookie: %w", err)
	}

	sessionJSON, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("invalid session cookie: %w", err)
	}

	var sessionData sessionCookie
	if err := json.Unmarshal(sessionJSON, &sessionData); err != nil {
		return nil, fmt.Errorf("invalid session data: %w", err)
	}

	user, err := h.store.GetUser(sessionData.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.Token.IsExpired() {
		newToken, err := h.auth.RefreshToken(user.Token.RefreshToken, user.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}
		user.Token = newToken
		h.store.SaveUser(user)
	}

	return user, nil
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
