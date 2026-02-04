// Spotify Web Playback SDK manager
let spotifyPlayer = null;
let deviceId = null;
let playerReady = false;
let accessToken = null;

// Define the SDK ready callback globally (must be before SDK loads)
window.onSpotifyWebPlaybackSDKReady = () => {
    console.log('Spotify SDK Ready - waiting for initialization...');
};

// Initialize Spotify Web Playback SDK
async function initializeSpotifyPlayer() {
    try {
        const tokenResponse = await getAccessToken();
        accessToken = tokenResponse.accessToken;

        window.onSpotifyWebPlaybackSDKReady = () => {
            console.log('Initializing Spotify Player...');
            const player = new Spotify.Player({
                name: 'Spotify Heardle',
                getOAuthToken: cb => { cb(accessToken); },
                volume: 0.5
            });

            player.addListener('ready', ({ device_id }) => {
                console.log('Player ready with Device ID:', device_id);
                deviceId = device_id;
                playerReady = true;
                updatePlayerStatus('Player ready!');
            });

            player.addListener('not_ready', ({ device_id }) => {
                console.log('Device ID has gone offline', device_id);
                playerReady = false;
                updatePlayerStatus('Player offline');
            });

            player.addListener('initialization_error', ({ message }) => {
                console.error('Initialization error:', message);
                updatePlayerStatus('Error: ' + message);
            });

            player.addListener('authentication_error', ({ message }) => {
                console.error('Authentication error:', message);
                updatePlayerStatus('Auth error: ' + message);
            });

            player.addListener('account_error', ({ message }) => {
                console.error('Account error:', message);
                updatePlayerStatus('Premium account required');
            });

            player.addListener('playback_error', ({ message }) => {
                console.error('Playback error:', message);
                updatePlayerStatus('Playback error: ' + message);
            });

            player.connect().then(success => {
                if (success) {
                    console.log('Spotify Player connected successfully');
                }
            });

            spotifyPlayer = player;
        };

        // If SDK is already loaded, trigger the callback
        if (window.Spotify) {
            window.onSpotifyWebPlaybackSDKReady();
        }
    } catch (error) {
        console.error('Failed to initialize player:', error);
        updatePlayerStatus('Failed to initialize player');
    }
}

// Play track with duration limit
async function playTrackWithLimit(trackUri, durationSeconds) {
    if (!playerReady || !deviceId) {
        throw new Error('Player not ready');
    }

    if (!accessToken) {
        const tokenResponse = await getAccessToken();
        accessToken = tokenResponse.accessToken;
    }

    const response = await fetch(`https://api.spotify.com/v1/me/player/play?device_id=${deviceId}`, {
        method: 'PUT',
        body: JSON.stringify({ uris: [trackUri] }),
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
    });

    if (!response.ok) {
        throw new Error(`Failed to play track: ${response.status}`);
    }

    setTimeout(async () => {
        await pausePlayback();
    }, durationSeconds * 1000);
}

// Pause playback
async function pausePlayback() {
    if (spotifyPlayer) {
        await spotifyPlayer.pause();
    }
}

// Update player status message
function updatePlayerStatus(message) {
    const statusElement = document.getElementById('player-status');
    if (statusElement) {
        statusElement.textContent = message;
        setTimeout(() => {
            statusElement.textContent = '';
        }, 3000);
    }
}
