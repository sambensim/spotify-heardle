// Package models defines data structures for the application.
package models

import "time"

// User represents a Spotify user.
type User struct {
	ID          string
	DisplayName string
	Token       *Token
}

// Token represents Spotify OAuth tokens.
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// IsExpired checks if the token is expired or will expire within 1 minute.
func (t *Token) IsExpired() bool {
	return time.Now().Add(1 * time.Minute).After(t.ExpiresAt)
}

// NewUser creates a new User instance.
func NewUser(id, displayName string, token *Token) *User {
	return &User{
		ID:          id,
		DisplayName: displayName,
		Token:       token,
	}
}
