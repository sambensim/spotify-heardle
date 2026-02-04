// Package models defines data structures for the application.
package models

const MaxGuesses = 3

// GameSession represents an active game session.
type GameSession struct {
	ID          string
	UserID      string
	PlaylistID  string
	CorrectSong Track
	Guesses     []Guess
	GuessesUsed int
	SkipsUsed   int
	IsComplete  bool
	Won         bool
}

// Track represents a Spotify track.
type Track struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []string `json:"artists"`
	PreviewURL string   `json:"previewUrl"`
}

// Guess represents a user's guess.
type Guess struct {
	TrackID   string
	TrackName string
	IsCorrect bool
}

// NewGameSession creates a new game session.
func NewGameSession(sessionID, userID, playlistID string, correctSong Track) *GameSession {
	return &GameSession{
		ID:          sessionID,
		UserID:      userID,
		PlaylistID:  playlistID,
		CorrectSong: correctSong,
		Guesses:     []Guess{},
		GuessesUsed: 0,
		SkipsUsed:   0,
		IsComplete:  false,
		Won:         false,
	}
}

// AddGuess adds a guess to the session and updates state.
func (s *GameSession) AddGuess(guess Guess) {
	s.Guesses = append(s.Guesses, guess)
	s.GuessesUsed++

	if guess.IsCorrect {
		s.IsComplete = true
		s.Won = true
	} else if s.GuessesUsed >= MaxGuesses {
		s.IsComplete = true
		s.Won = false
	}
}

// GetAudioDuration returns the audio duration in seconds based on guesses and skips used.

// The durations array contains cumulative clip lengths: [1, 3, 6, 10, 15] seconds.
// Each increment adds 1 more second than the previous: 1s, then +2s, then +3s, then +4s, then +5s.
func (s *GameSession) GetAudioDuration() int {
	durations := []int{1, 3, 6, 10, 15}
	totalSteps := s.GuessesUsed + s.SkipsUsed
	if totalSteps >= len(durations) {
		// After exhausting the array, continue adding 5 seconds per step
		excess := totalSteps - len(durations) + 1
		return durations[len(durations)-1] + (excess * 5)
	}
	return durations[totalSteps]
}

// GetTotalAudioDuration returns the cumulative audio duration revealed so far.
// Returns the cumulative duration from the durations array based on steps completed.
func (s *GameSession) GetTotalAudioDuration() int {
	durations := []int{1, 3, 6, 10, 15}
	totalSteps := s.GuessesUsed + s.SkipsUsed
	if totalSteps == 0 {
		return 0
	}
	if totalSteps >= len(durations) {
		// After exhausting the array, continue adding the increment (5 seconds)
		excess := totalSteps - len(durations)
		return durations[len(durations)-1] + (excess * 5)
	}
	return durations[totalSteps-1]
}

// GetNextAudioDuration returns what the next audio duration would be.
// Returns the cumulative duration at the next step.
func (s *GameSession) GetNextAudioDuration() int {
	durations := []int{1, 3, 6, 10, 15}
	nextStep := s.GuessesUsed + s.SkipsUsed
	if nextStep >= len(durations) {
		// After exhausting the array, continue adding 5 seconds per step
		excess := nextStep - len(durations) + 1
		return durations[len(durations)-1] + (excess * 5)
	}
	return durations[nextStep]
}

// CanSkip returns true if the user can still skip (next skip won't exceed 20 seconds total).
func (s *GameSession) CanSkip() bool {
	if s.IsComplete {
		return false
	}
	nextTotal := s.GetTotalAudioDuration() + s.GetNextAudioDuration()
	return nextTotal <= 20
}

// Skip increments the skip counter.
func (s *GameSession) Skip() {
	s.SkipsUsed++
}

// MarkComplete marks the session as complete.
func (s *GameSession) MarkComplete(won bool) {
	s.IsComplete = true
	s.Won = won
}
