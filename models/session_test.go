// Package models defines data structures for the application.
package models

import (
	"testing"
)

func TestNewGameSession(t *testing.T) {
	sessionID := "session123"
	userID := "user456"
	playlistID := "playlist789"
	correctSong := Track{
		ID:         "track1",
		Name:       "Song Name",
		Artists:    []string{"Artist 1"},
		PreviewURL: "http://preview.url",
	}

	session := NewGameSession(sessionID, userID, playlistID, correctSong)

	if session.ID != sessionID {
		t.Errorf("ID = %q, want %q", session.ID, sessionID)
	}

	if session.UserID != userID {
		t.Errorf("UserID = %q, want %q", session.UserID, userID)
	}

	if session.PlaylistID != playlistID {
		t.Errorf("PlaylistID = %q, want %q", session.PlaylistID, playlistID)
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
		PlaylistID:  "playlist1",
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
		PlaylistID:  "playlist1",
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
		PlaylistID:  "playlist1",
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
		skipsUsed   int
		want        int
	}{
		{"first attempt", 0, 0, 1},
		{"after one guess", 1, 0, 3},
		{"after two guesses", 2, 0, 6},
		{"after one skip", 0, 1, 3},
		{"after one guess and one skip", 1, 1, 6},
		{"after two skips", 0, 2, 6},
		{"max duration", 5, 0, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &GameSession{
				GuessesUsed: tt.guessesUsed,
				SkipsUsed:   tt.skipsUsed,
			}
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

func TestGameSessionGetTotalAudioDuration(t *testing.T) {
	tests := []struct {
		name        string
		guessesUsed int
		skipsUsed   int
		want        int
	}{
		{"no guesses or skips", 0, 0, 0},
		{"one guess", 1, 0, 1},
		{"two guesses", 2, 0, 3},   // cumulative: 1s + 2s more
		{"one skip", 0, 1, 1},
		{"one guess and one skip", 1, 1, 3}, // cumulative: 1s + 2s more
		{"three total", 2, 1, 6},   // cumulative: 1s + 2s + 3s more
		{"five total", 3, 2, 15},   // cumulative: 1s + 2s + 3s + 4s + 5s more
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &GameSession{
				GuessesUsed: tt.guessesUsed,
				SkipsUsed:   tt.skipsUsed,
			}
			got := session.GetTotalAudioDuration()
			if got != tt.want {
				t.Errorf("GetTotalAudioDuration() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGameSessionGetNextAudioDuration(t *testing.T) {
	tests := []struct {
		name        string
		guessesUsed int
		skipsUsed   int
		want        int
	}{
		{"first attempt", 0, 0, 1},
		{"second attempt", 1, 0, 3},
		{"third attempt", 0, 2, 6},
		{"fourth attempt", 2, 1, 10},
		{"fifth attempt", 3, 1, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &GameSession{
				GuessesUsed: tt.guessesUsed,
				SkipsUsed:   tt.skipsUsed,
			}
			got := session.GetNextAudioDuration()
			if got != tt.want {
				t.Errorf("GetNextAudioDuration() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGameSessionCanSkip(t *testing.T) {
	tests := []struct {
		name        string
		guessesUsed int
		skipsUsed   int
		isComplete  bool
		want        bool
	}{
		{"can skip at start", 0, 0, false, true},
		{"can skip after one", 1, 0, false, true},
		{"can skip after two", 0, 2, false, true},
		{"can skip at 15 seconds", 2, 2, false, true},     // cumulative: 10s revealed, next is 15s total
		{"can skip at 15 seconds total", 3, 1, false, true},  // cumulative: 10s revealed, next is 15s total
		{"can skip at 20 seconds", 4, 1, false, true}, // cumulative: 15s revealed, next would be 20s total (15s + 5s more)
		{"cannot skip at 65 seconds", 0, 14, false, false}, // cumulative: 60s revealed, next would be 65s total (60s + 5s more)
		{"cannot skip when complete", 1, 1, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &GameSession{
				GuessesUsed: tt.guessesUsed,
				SkipsUsed:   tt.skipsUsed,
				IsComplete:  tt.isComplete,
			}
			got := session.CanSkip()
			if got != tt.want {
				t.Errorf("CanSkip() = %v, want %v (total=%d, next=%d)", 
					got, tt.want, session.GetTotalAudioDuration(), session.GetNextAudioDuration())
			}
		})
	}
}

func TestGameSessionSkip(t *testing.T) {
	session := &GameSession{
		SkipsUsed: 0,
	}

	session.Skip()

	if session.SkipsUsed != 1 {
		t.Errorf("SkipsUsed = %d, want 1", session.SkipsUsed)
	}
}
