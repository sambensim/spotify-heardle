// Package spotify provides Spotify API client functionality.
package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"spotify-heardle/models"
)

const apiBaseURL = "https://api.spotify.com/v1"

// Client is a Spotify API client.
type Client struct {
	token *models.Token
}

// Playlist represents a Spotify playlist.
type Playlist struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Images []PlaylistImage `json:"images"`
	Tracks TracksInfo     `json:"tracks"`
}

// PlaylistImage represents a playlist cover image.
type PlaylistImage struct {
	URL string `json:"url"`
}

// TracksInfo contains playlist track metadata.
type TracksInfo struct {
	Total int `json:"total"`
}

// UserProfile represents a Spotify user profile.
type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type playlistsResponse struct {
	Items []Playlist `json:"items"`
}

type playlistTracksResponse struct {
	Items []struct {
		Track trackInfo `json:"track"`
	} `json:"items"`
}

type trackInfo struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []artist `json:"artists"`
	PreviewURL string   `json:"preview_url"`
}

type artist struct {
	Name string `json:"name"`
}

type searchResponse struct {
	Tracks struct {
		Items []trackInfo `json:"items"`
	} `json:"tracks"`
}

// NewClient creates a new Spotify API client.
func NewClient(token *models.Token) *Client {
	return &Client{token: token}
}

// GetUserProfile retrieves the current user's profile.
func (c *Client) GetUserProfile() (*UserProfile, error) {
	endpoint := apiBaseURL + "/me"
	
	var profile UserProfile
	if err := c.makeRequest("GET", endpoint, &profile); err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}

	return &profile, nil
}

// GetUserPlaylists retrieves the current user's playlists.
func (c *Client) GetUserPlaylists() ([]Playlist, error) {
	endpoint := apiBaseURL + "/me/playlists?limit=50"
	
	var response playlistsResponse
	if err := c.makeRequest("GET", endpoint, &response); err != nil {
		return nil, fmt.Errorf("getting playlists: %w", err)
	}

	return response.Items, nil
}

// GetPlaylistTracks retrieves tracks from a playlist.
func (c *Client) GetPlaylistTracks(playlistID string) ([]models.Track, error) {
	endpoint := fmt.Sprintf("%s/playlists/%s/tracks?limit=50", apiBaseURL, playlistID)
	
	var response playlistTracksResponse
	if err := c.makeRequest("GET", endpoint, &response); err != nil {
		return nil, fmt.Errorf("getting playlist tracks: %w", err)
	}

	tracks := make([]models.Track, 0, len(response.Items))
	for _, item := range response.Items {
		artists := make([]string, len(item.Track.Artists))
		for i, artist := range item.Track.Artists {
			artists[i] = artist.Name
		}

		tracks = append(tracks, models.Track{
			ID:         item.Track.ID,
			Name:       item.Track.Name,
			Artists:    artists,
			PreviewURL: item.Track.PreviewURL,
		})
	}

	return tracks, nil
}

// SearchTracks searches for tracks by query.
func (c *Client) SearchTracks(query string) ([]models.Track, error) {
	endpoint := fmt.Sprintf("%s/search?q=%s&type=track&limit=20", apiBaseURL, url.QueryEscape(query))
	
	var response searchResponse
	if err := c.makeRequest("GET", endpoint, &response); err != nil {
		return nil, fmt.Errorf("searching tracks: %w", err)
	}

	tracks := make([]models.Track, 0, len(response.Tracks.Items))
	for _, item := range response.Tracks.Items {
		artists := make([]string, len(item.Artists))
		for i, artist := range item.Artists {
			artists[i] = artist.Name
		}

		tracks = append(tracks, models.Track{
			ID:         item.ID,
			Name:       item.Name,
			Artists:    artists,
			PreviewURL: item.PreviewURL,
		})
	}

	return tracks, nil
}

func (c *Client) makeRequest(method, endpoint string, result interface{}) error {
	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}
