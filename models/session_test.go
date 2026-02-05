// Package models defines data structures for the application.
package models

import (
	"testing"
)

func TestNewGameSession(t *testing.T) {
	sessionID := "session123"
	userID := "user456"
	playlistIDs := []string{"playlist789"}
	correctSong := Track{
		ID:         "track1",
		Name:       "Song Name",
		Artists:    []string{"Artist 1"},
		PreviewURL: "http://preview.url",
	}

	session := NewGameSession(sessionID, userID, playlistIDs, correctSong)

	if session.ID != sessionID {
		t.Errorf("ID = %q, want %q", session.ID, sessionID)
	}

	if session.UserID != userID {
		t.Errorf("UserID = %q, want %q", session.UserID, userID)
	}

	if len(session.PlaylistIDs) != 1 || session.PlaylistIDs[0] != playlistIDs[0] {
		t.Errorf("PlaylistIDs = %v, want %v", session.PlaylistIDs, playlistIDs)
	}

	if session.CorrectSong.ID != correctSong.ID {
		t.Errorf("CorrectSong.ID = %q, want %q", session.CorrectSong.ID, correctSong.ID)
	}

	if session.GuessesUsed != 0 {
		t.Errorf("GuessesUsed = %d, want 0", session.GuessesUsed)
	}

	if session.IsComplete {
		t.Error("IsComplete = true, want false")
	}

	if session.Won {
		t.Error("Won = true, want false")
	}

	if len(session.Guesses) != 0 {
		t.Errorf("len(Guesses) = %d, want 0", len(session.Guesses))
	}
}

func TestGameSessionAddGuess(t *testing.T) {
	session := &GameSession{
		ID:          "session1",
		UserID:      "user1",
		PlaylistIDs: []string{"playlist1"},
		CorrectSong: Track{ID: "correct_track"},
		GuessesUsed: 0,
		IsComplete:  false,
	}

	guess := Guess{
		TrackID:   "track1",
		TrackName: "Wrong Song",
		IsCorrect: false,
	}

	session.AddGuess(guess)

	if session.GuessesUsed != 1 {
		t.Errorf("GuessesUsed = %d, want 1", session.GuessesUsed)
	}

	if len(session.Guesses) != 1 {
		t.Fatalf("len(Guesses) = %d, want 1", len(session.Guesses))
	}

	if session.Guesses[0].TrackID != guess.TrackID {
		t.Errorf("Guesses[0].TrackID = %q, want %q", session.Guesses[0].TrackID, guess.TrackID)
	}
}

func TestGameSessionAddCorrectGuess(t *testing.T) {
	session := &GameSession{
		ID:          "session1",
		UserID:      "user1",
		PlaylistIDs: []string{"playlist1"},
		CorrectSong: Track{ID: "correct_track"},
		GuessesUsed: 0,
		IsComplete:  false,
		Won:         false,
	}

	correctGuess := Guess{
		TrackID:   "correct_track",
		TrackName: "Correct Song",
		IsCorrect: true,
	}

	session.AddGuess(correctGuess)

	if !session.IsComplete {
		t.Error("IsComplete = false, want true after correct guess")
	}

	if !session.Won {
		t.Error("Won = false, want true after correct guess")
	}
}

func TestGameSessionAddGuessReachesMax(t *testing.T) {
	session := &GameSession{
		ID:          "session1",
		UserID:      "user1",
		PlaylistIDs: []string{"playlist1"},
		CorrectSong: Track{ID: "correct_track"},
		GuessesUsed: 2,
		IsComplete:  false,
		Won:         false,
		Guesses: []Guess{
			{TrackID: "wrong1", IsCorrect: false},
			{TrackID: "wrong2", IsCorrect: false},
		},
	}

	wrongGuess := Guess{
		TrackID:   "wrong3",
		TrackName: "Wrong Song",
		IsCorrect: false,
	}

	session.AddGuess(wrongGuess)

	if session.GuessesUsed != 3 {
		t.Errorf("GuessesUsed = %d, want 3", session.GuessesUsed)
	}

	if !session.IsComplete {
		t.Error("IsComplete = false, want true after 3 guesses")
	}

	if session.Won {
		t.Error("Won = true, want false after 3 wrong guesses")
	}
}

func TestGameSessionGetAudioDuration(t *testing.T) {
	tests := []struct {
		name        string
		guessesUsed int
		want        int
	}{
		{"first guess", 0, 1},
		{"second guess", 1, 2},
		{"third guess", 2, 4},
		{"beyond max", 3, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &GameSession{GuessesUsed: tt.guessesUsed}
			got := session.GetAudioDuration()
			if got != tt.want {
				t.Errorf("GetAudioDuration() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGameSessionMarkComplete(t *testing.T) {
	session := &GameSession{
		ID:         "session1",
		IsComplete: false,
	}

	session.MarkComplete(true)

	if !session.IsComplete {
		t.Error("IsComplete = false, want true")
	}

	if !session.Won {
		t.Error("Won = false, want true")
	}
}
