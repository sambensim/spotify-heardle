// Package handlers provides HTTP request handlers.
package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"spotify-heardle/models"
	"spotify-heardle/spotify"
	"spotify-heardle/storage"
	"time"
)

// GameHandler handles game-related routes.
type GameHandler struct {
	auth  *AuthHandler
	store *storage.MemoryStore
}

type startGameRequest struct {
	PlaylistID string `json:"playlistId"`
}

type startGameResponse struct {
	SessionID     string `json:"sessionId"`
	AudioDuration int    `json:"audioDuration"`
	PreviewURL    string `json:"previewUrl"`
}

type submitGuessRequest struct {
	SessionID string `json:"sessionId"`
	TrackID   string `json:"trackId"`
	TrackName string `json:"trackName"`
}

type submitGuessResponse struct {
	IsCorrect     bool          `json:"isCorrect"`
	IsComplete    bool          `json:"isComplete"`
	Won           bool          `json:"won"`
	GuessesUsed   int           `json:"guessesUsed"`
	AudioDuration int           `json:"audioDuration"`
	CorrectSong   *models.Track `json:"correctSong,omitempty"`
}

type skipRequest struct {
	SessionID string `json:"sessionId"`
}

type skipResponse struct {
	CorrectSong models.Track `json:"correctSong"`
}

// NewGameHandler creates a new game handler.
func NewGameHandler(auth *AuthHandler, store *storage.MemoryStore) *GameHandler {
	rand.Seed(time.Now().UnixNano())
	return &GameHandler{
		auth:  auth,
		store: store,
	}
}

// HandleStartGame starts a new game session.
func (h *GameHandler) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req startGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	client := spotify.NewClient(user.Token)
	tracks, err := client.GetPlaylistTracks(req.PlaylistID)
	if err != nil {
		http.Error(w, "Failed to get playlist tracks", http.StatusInternalServerError)
		return
	}

	tracksWithPreview := filterTracksWithPreview(tracks)
	if len(tracksWithPreview) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("No preview URLs found in the first 100 tracks. Spotify doesn't provide previews for all songs. Try a different playlist with more mainstream/popular tracks."),
		})
		return
	}

	selectedTrack := selectRandomTrack(tracksWithPreview)

	sessionID, err := generateSessionID()
	if err != nil {
		http.Error(w, "Failed to generate session ID", http.StatusInternalServerError)
		return
	}

	session := models.NewGameSession(sessionID, user.ID, req.PlaylistID, selectedTrack)
	if err := h.store.SaveSession(session); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	response := startGameResponse{
		SessionID:     sessionID,
		AudioDuration: session.GetAudioDuration(),
		PreviewURL:    selectedTrack.PreviewURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleSubmitGuess processes a user's guess.
func (h *GameHandler) HandleSubmitGuess(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req submitGuessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := h.store.GetSession(req.SessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	if session.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if session.IsComplete {
		http.Error(w, "Game already complete", http.StatusBadRequest)
		return
	}

	isCorrect := req.TrackID == session.CorrectSong.ID
	guess := models.Guess{
		TrackID:   req.TrackID,
		TrackName: req.TrackName,
		IsCorrect: isCorrect,
	}

	session.AddGuess(guess)
	h.store.SaveSession(session)

	response := submitGuessResponse{
		IsCorrect:     isCorrect,
		IsComplete:    session.IsComplete,
		Won:           session.Won,
		GuessesUsed:   session.GuessesUsed,
		AudioDuration: session.GetAudioDuration(),
	}

	if session.IsComplete {
		response.CorrectSong = &session.CorrectSong
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleSkip skips the current game and reveals the answer.
func (h *GameHandler) HandleSkip(w http.ResponseWriter, r *http.Request) {
	user, err := h.auth.GetUserFromSession(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req skipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := h.store.GetSession(req.SessionID)
	if err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	if session.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	session.MarkComplete(false)
	h.store.SaveSession(session)

	response := skipResponse{
		CorrectSong: session.CorrectSong,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func filterTracksWithPreview(tracks []models.Track) []models.Track {
	filtered := make([]models.Track, 0)
	for _, track := range tracks {
		if track.PreviewURL != "" {
			filtered = append(filtered, track)
		}
	}
	return filtered
}

func selectRandomTrack(tracks []models.Track) models.Track {
	if len(tracks) == 0 {
		return models.Track{}
	}
	return tracks[rand.Intn(len(tracks))]
}
