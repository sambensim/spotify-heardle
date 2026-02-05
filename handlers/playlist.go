// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"net/http"
	"spotify-heardle/spotify"
)

// PlaylistHandler handles playlist-related routes.
type PlaylistHandler struct {
	auth *AuthHandler
}

// NewPlaylistHandler creates a new playlist handler.
func NewPlaylistHandler(auth *AuthHandler) *PlaylistHandler {
	return &PlaylistHandler{auth: auth}
}

// HandleGetPlaylists returns user's playlists.
func (h *PlaylistHandler) HandleGetPlaylists(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	client := spotify.NewClient(user.Token)
	playlists, err := client.GetUserPlaylists()
	if err != nil {
		http.Error(w, "Failed to get playlists", http.StatusInternalServerError)
		return
	}

	// Add "Liked Songs" as a special playlist option
	likedSongsPlaylist := spotify.Playlist{
		ID:     "liked_songs",
		Name:   "Liked Songs",
		Images: []spotify.PlaylistImage{},
		Tracks: spotify.TracksInfo{Total: 0}, // Total is set to 0 to avoid an additional API call
	}

	// Prepend liked songs to the list
	allPlaylists := append([]spotify.Playlist{likedSongsPlaylist}, playlists...)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allPlaylists)
}
