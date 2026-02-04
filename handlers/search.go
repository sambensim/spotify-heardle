// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
	"spotify-heardle/spotify"
)

// SearchHandler handles track search routes.
type SearchHandler struct {
	auth *AuthHandler
}

// NewSearchHandler creates a new search handler.
func NewSearchHandler(auth *AuthHandler) *SearchHandler {
	return &SearchHandler{auth: auth}
}

// HandleSearch searches for tracks.
func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	client := spotify.NewClient(user.Token)
	tracks, err := client.SearchTracks(query)
	if err != nil {
		http.Error(w, "Failed to search tracks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}
