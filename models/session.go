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
	IsComplete  bool
	Won         bool
}

// Track represents a Spotify track.
type Track struct {
	ID         string
	Name       string
	Artists    []string
	PreviewURL string
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

// GetAudioDuration returns the audio duration in seconds based on guesses used.
func (s *GameSession) GetAudioDuration() int {
	durations := []int{1, 2, 4}
	if s.GuessesUsed >= len(durations) {
		return durations[len(durations)-1]
	}
	return durations[s.GuessesUsed]
}

// MarkComplete marks the session as complete.
func (s *GameSession) MarkComplete(won bool) {
	s.IsComplete = true
	s.Won = won
}
