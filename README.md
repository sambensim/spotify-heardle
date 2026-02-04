# Spotify Heardle Clone

A web-based music guessing game inspired by Heardle, where users can play with their own Spotify playlists.

## Features

- Spotify OAuth authentication
- Choose from your own playlists
- **Full track playback** using Spotify Web Playback SDK
- Progressive audio reveal (1s → 3s → 6s → 10s → 15s)
- Skip to hear more without guessing
- 3 guesses per game
- Search Spotify tracks to make guesses
- Unlimited plays with the same playlist

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
- **Audio Duration**: Progressively reveals 1s → 3s → 6s → 10s → 15s clips
- **Skip Limit**: Can skip until total revealed audio would exceed 60 seconds

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
