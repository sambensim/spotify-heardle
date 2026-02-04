// API client for backend communication
const API_BASE = '';

async function fetchAPI(endpoint, options = {}) {
    const response = await fetch(API_BASE + endpoint, {
        ...options,
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    });

    if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
    }

    return response.json();
}

async function getPlaylists() {
    return fetchAPI('/api/playlists');
}

async function searchTracks(query) {
    return fetchAPI(`/api/search?q=${encodeURIComponent(query)}`);
}

async function startGame(playlistId) {
    return fetchAPI('/api/game/start', {
        method: 'POST',
        body: JSON.stringify({ playlistId }),
    });
}

async function submitGuess(sessionId, trackId, trackName) {
    return fetchAPI('/api/game/guess', {
        method: 'POST',
        body: JSON.stringify({ sessionId, trackId, trackName }),
    });
}

async function skipCurrentGame(sessionId) {
    return fetchAPI('/api/game/skip', {
        method: 'POST',
        body: JSON.stringify({ sessionId }),
    });
}

async function logout() {
    return fetchAPI('/api/logout', { method: 'POST' });
}
