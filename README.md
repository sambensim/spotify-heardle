# Spotify Heardle Clone

A web-based music guessing game inspired by Heardle, where users can play with their own Spotify playlists.

## Features

- Spotify OAuth authentication
- **Multiple playlist selection** - Choose one or more playlists to play from
- **Liked Songs support** - Play using your saved/liked tracks
- Songs are randomly selected from the combined pool of all selected playlists
- **Full track playback** using Spotify Web Playback SDK
- Progressive audio reveal (1s → 2s → 4s)
- 3 guesses per game
- Search Spotify tracks to make guesses
- Unlimited plays

## Prerequisites

- Go 1.21 or higher
- **Spotify Premium Account** (required for Web Playback SDK)
- Spotify Developer Account

## Setup

1. **Get Spotify API Credentials**
   - Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
   - Create a new app
   - Add `http://localhost:8080/callback` to Redirect URIs
   - Note your Client ID and Client Secret

2. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your Spotify credentials
   ```

3. **Run the Application**
   ```bash
   ./run.sh
   ```
   
   Or manually with:
   ```bash
   source .env
   go run main.go
   ```

4. **Open Browser**
   Navigate to `http://localhost:8080`

## Development

### Run Tests
```bash
go test ./...
```

### Run Single Test
```bash
go test ./path/to/package -run TestName
```

## Project Structure

```
├── main.go              # Entry point
├── config/              # Configuration management
├── handlers/            # HTTP handlers
├── models/              # Data models
├── spotify/             # Spotify API client
├── storage/             # Session storage
└── static/              # Frontend assets
    ├── css/
    ├── js/
    │   ├── api.js       # Backend API client
    │   ├── player.js    # Spotify Web Playback SDK
    │   ├── game.js      # Game logic
    │   └── ...
    └── *.html
```

## Technical Notes

- **Web Playback SDK**: Requires Spotify Premium and uses the browser-based player
- **Authentication**: OAuth 2.0 with PKCE flow
- **Storage**: In-memory (sessions cleared on restart)
- **Audio Duration**: Progressively reveals 1s → 2s → 4s clips

## Troubleshooting

**"Player not ready"**
- Ensure you have Spotify Premium
- Close any other Spotify players (app, web)
- Refresh the page

**"Account error"**
- Spotify Premium is required for Web Playback SDK
- Verify your account at https://www.spotify.com/account

**"Authentication error"**
- Log out and log in again
- Check that redirect URI matches in Spotify Dashboard

## License

MIT
