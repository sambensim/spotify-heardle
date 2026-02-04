// Game logic and state management
let gameState = {
    sessionId: null,
    guessesUsed: 0,
    audioDuration: 1,
    trackUri: null,
    isComplete: false,
};

document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const playlistsParam = urlParams.get('playlists');
    const legacyPlaylistId = urlParams.get('playlist');

    let playlistIds = [];
    
    // Support new multi-playlist format
    if (playlistsParam) {
        try {
            playlistIds = JSON.parse(decodeURIComponent(playlistsParam));
        } catch (e) {
            showError('Invalid playlist parameter');
            return;
        }
    } 
    // Backward compatibility with old single playlist format
    else if (legacyPlaylistId) {
        playlistIds = [legacyPlaylistId];
    } 
    else {
        showError('No playlist selected');
        return;
    }

    if (!Array.isArray(playlistIds) || playlistIds.length === 0) {
        showError('No playlist selected');
        return;
    }

    showLoadingMessage('Initializing Spotify player...');
    await initializeSpotifyPlayer();
    
    showLoadingMessage('Starting game...');
    await initializeGame(playlistIds);
    initSearch();
});

function showLoadingMessage(message) {
    const loading = document.getElementById('loading');
    if (loading) {
        loading.textContent = message;
        loading.style.display = 'block';
    }
}

async function initializeGame(playlistIds) {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const gameContainer = document.getElementById('game-container');

    try {
        const response = await startGame(playlistIds);
        
        gameState.sessionId = response.sessionId;
        gameState.audioDuration = response.audioDuration;
        gameState.trackUri = response.trackUri;

        updateGameUI();

        loading.style.display = 'none';
        gameContainer.style.display = 'block';
    } catch (err) {
        loading.style.display = 'none';
        error.textContent = 'Failed to start game. Please try a different playlist or check your Spotify Premium subscription.';
        error.style.display = 'block';
        console.error('Error starting game:', err);
    }
}

async function playAudio() {
    if (!playerReady) {
        showError('Player not ready. Please wait...');
        return;
    }

    const playBtn = document.getElementById('play-btn');
    playBtn.disabled = true;
    playBtn.textContent = '▶ Playing...';

    try {
        await playTrackWithLimit(gameState.trackUri, gameState.audioDuration);
        
        setTimeout(() => {
            playBtn.disabled = false;
            playBtn.textContent = '▶ Play Again';
        }, gameState.audioDuration * 1000 + 500);
    } catch (error) {
        console.error('Playback failed:', error);
        showError('Playback failed. Make sure Spotify is not playing elsewhere.');
        playBtn.disabled = false;
        playBtn.textContent = '▶ Play';
    }
}

async function handleGuess(trackId, trackName) {
    if (gameState.isComplete) {
        return;
    }

    const searchInput = document.getElementById('search-input');
    searchInput.value = '';
    searchInput.disabled = true;

    try {
        const response = await submitGuess(gameState.sessionId, trackId, trackName);
        
        gameState.guessesUsed = response.guessesUsed;
        gameState.audioDuration = response.audioDuration;
        gameState.isComplete = response.isComplete;

        addGuessToList(trackName, response.isCorrect);
        updateGameUI();

        if (response.isComplete) {
            showResult(response.won, response.correctSong);
        } else {
            searchInput.disabled = false;
            searchInput.focus();
        }
    } catch (error) {
        console.error('Guess submission failed:', error);
        searchInput.disabled = false;
        showError('Failed to submit guess');
    }
}

async function skipGame() {
    if (!confirm('Are you sure you want to skip and see the answer?')) {
        return;
    }

    try {
        const response = await skipCurrentGame(gameState.sessionId);
        gameState.isComplete = true;
        showResult(false, response.correctSong);
    } catch (error) {
        console.error('Skip failed:', error);
        showError('Failed to skip');
    }
}

function updateGameUI() {
    document.getElementById('guesses-used').textContent = gameState.guessesUsed;
    document.getElementById('audio-duration').textContent = gameState.audioDuration;

    const skipBtn = document.getElementById('skip-btn');
    const searchSection = document.getElementById('search-section');
    
    if (gameState.isComplete) {
        skipBtn.style.display = 'none';
        searchSection.style.display = 'none';
    }
}

function addGuessToList(trackName, isCorrect) {
    const guessesList = document.getElementById('guesses-list');
    const guessItem = document.createElement('div');
    guessItem.className = `guess-item ${isCorrect ? 'guess-correct' : 'guess-incorrect'}`;
    guessItem.innerHTML = `
        <span>${trackName}</span>
        <span>${isCorrect ? '✓ Correct!' : '✗ Incorrect'}</span>
    `;
    guessesList.appendChild(guessItem);
}

function showResult(won, correctSong) {
    const modal = document.getElementById('result-modal');
    const title = document.getElementById('result-title');
    const songDiv = document.getElementById('result-song');

    title.textContent = won ? '🎉 You Win!' : '😔 Game Over';
    
    const artists = correctSong.artists.join(', ');
    songDiv.innerHTML = `
        <div class="result-song-name">${correctSong.name}</div>
        <div class="result-song-artist">${artists}</div>
    `;

    modal.style.display = 'flex';
}

function showError(message) {
    const error = document.getElementById('error');
    error.textContent = message;
    error.style.display = 'block';
}

function newGame() {
    window.location.href = '/playlists.html';
}
