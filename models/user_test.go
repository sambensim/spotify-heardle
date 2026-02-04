// Package models defines data structures for the application.
package models

import (
	"testing"
	"time"
)

func TestTokenIsExpired(t *testing.T) {
	tests := []struct {
		name    string
		token   *Token
		want    bool
	}{
		{
			name: "expired token",
			token: &Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(-1 * time.Hour),
			},
			want: true,
		},
		{
			name: "valid token",
			token: &Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "token expiring soon",
			token: &Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(30 * time.Second),
			},
			want: true,
		},
		{
			name: "token valid with buffer",
			token: &Token{
				AccessToken:  "access",
				RefreshToken: "refresh",
				ExpiresAt:    time.Now().Add(2 * time.Minute),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.token.IsExpired()
			if got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	userID := "user123"
	displayName := "Test User"
	token := &Token{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		ExpiresAt:    time.Now().Add(1 * time.Hour),
	}

	user := NewUser(userID, displayName, token)

	if user.ID != userID {
		t.Errorf("ID = %q, want %q", user.ID, userID)
	}

	if user.DisplayName != displayName {
		t.Errorf("DisplayName = %q, want %q", user.DisplayName, displayName)
	}

	if user.Token != token {
		t.Errorf("Token mismatch")
	}
}
