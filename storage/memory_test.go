// Package storage provides in-memory storage for sessions and users.
package storage

import (
	"spotify-heardle/models"
	"testing"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	if store == nil {
		t.Fatal("NewMemoryStore() returned nil")
	}
}

func TestSaveAndGetUser(t *testing.T) {
	store := NewMemoryStore()
	user := &models.User{
		ID:          "user123",
		DisplayName: "Test User",
	}

	err := store.SaveUser(user)
	if err != nil {
		t.Fatalf("SaveUser() failed: %v", err)
	}

	retrieved, err := store.GetUser("user123")
	if err != nil {
		t.Fatalf("GetUser() failed: %v", err)
	}

	if retrieved.ID != user.ID {
		t.Errorf("retrieved.ID = %q, want %q", retrieved.ID, user.ID)
	}

	if retrieved.DisplayName != user.DisplayName {
		t.Errorf("retrieved.DisplayName = %q, want %q", retrieved.DisplayName, user.DisplayName)
	}
}

func TestGetUserNotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.GetUser("nonexistent")
	if err == nil {
		t.Error("GetUser() succeeded for nonexistent user, want error")
	}
}

func TestSaveAndGetSession(t *testing.T) {
	store := NewMemoryStore()
	session := &models.GameSession{
		ID:           "session123",
		UserID:       "user456",
		PlaylistIDs:  []string{"playlist789"},
		CorrectSong:  models.Track{ID: "track1"},
		GuessesUsed:  0,
	}

	err := store.SaveSession(session)
	if err != nil {
		t.Fatalf("SaveSession() failed: %v", err)
	}

	retrieved, err := store.GetSession("session123")
	if err != nil {
		t.Fatalf("GetSession() failed: %v", err)
	}

	if retrieved.ID != session.ID {
		t.Errorf("retrieved.ID = %q, want %q", retrieved.ID, session.ID)
	}

	if retrieved.UserID != session.UserID {
		t.Errorf("retrieved.UserID = %q, want %q", retrieved.UserID, session.UserID)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	store := NewMemoryStore()

	_, err := store.GetSession("nonexistent")
	if err == nil {
		t.Error("GetSession() succeeded for nonexistent session, want error")
	}
}

func TestDeleteSession(t *testing.T) {
	store := NewMemoryStore()
	session := &models.GameSession{
		ID:     "session123",
		UserID: "user456",
	}

	store.SaveSession(session)

	err := store.DeleteSession("session123")
	if err != nil {
		t.Fatalf("DeleteSession() failed: %v", err)
	}

	_, err = store.GetSession("session123")
	if err == nil {
		t.Error("GetSession() succeeded after deletion, want error")
	}
}

func TestDeleteSessionNotFound(t *testing.T) {
	store := NewMemoryStore()

	err := store.DeleteSession("nonexistent")
	if err == nil {
		t.Error("DeleteSession() succeeded for nonexistent session, want error")
	}
}

func TestSaveUserOverwrites(t *testing.T) {
	store := NewMemoryStore()
	user1 := &models.User{
		ID:          "user123",
		DisplayName: "First Name",
	}
	user2 := &models.User{
		ID:          "user123",
		DisplayName: "Second Name",
	}

	store.SaveUser(user1)
	store.SaveUser(user2)

	retrieved, err := store.GetUser("user123")
	if err != nil {
		t.Fatalf("GetUser() failed: %v", err)
	}

	if retrieved.DisplayName != "Second Name" {
		t.Errorf("DisplayName = %q, want %q", retrieved.DisplayName, "Second Name")
	}
}
