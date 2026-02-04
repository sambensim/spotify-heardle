// Package storage provides in-memory storage for sessions and users.
package storage

import (
	"fmt"
	"spotify-heardle/models"
	"sync"
)

// MemoryStore implements in-memory storage.
type MemoryStore struct {
	users    map[string]*models.User
	sessions map[string]*models.GameSession
	mu       sync.RWMutex
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:    make(map[string]*models.User),
		sessions: make(map[string]*models.GameSession),
	}
}

// SaveUser stores a user.
func (s *MemoryStore) SaveUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	return nil
}

// GetUser retrieves a user by ID.
func (s *MemoryStore) GetUser(userID string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	return user, nil
}

// SaveSession stores a game session.
func (s *MemoryStore) SaveSession(session *models.GameSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	return nil
}

// GetSession retrieves a game session by ID.
func (s *MemoryStore) GetSession(sessionID string) (*models.GameSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}
	return session, nil
}

// DeleteSession removes a game session.
func (s *MemoryStore) DeleteSession(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[sessionID]; !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}
	delete(s.sessions, sessionID)
	return nil
}
