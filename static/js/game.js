// Game logic and state management
let gameState = {
    sessionId: null,
    guessesUsed: 0,
    audioDuration: 1,
    previewUrl: null,
    isComplete: false,
};

document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const playlistId = urlParams.get('playlist');

    if (!playlistId) {
        showError('No playlist selected');
        return;
    }

    await initializeGame(playlistId);
    initSearch();
});

async function initializeGame(playlistId) {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const gameContainer = document.getElementById('game-container');

    try {
        const response = await startGame(playlistId);
        
        gameState.sessionId = response.sessionId;
        gameState.audioDuration = response.audioDuration;
        gameState.previewUrl = response.previewUrl;

        const audio = document.getElementById('audio');
        audio.src = response.previewUrl;

        updateGameUI();

        loading.style.display = 'none';
        gameContainer.style.display = 'block';
    } catch (err) {
        loading.style.display = 'none';
        error.textContent = 'Failed to start game. The playlist may not have any tracks with previews.';
        error.style.display = 'block';
        console.error('Error starting game:', err);
    }
}

function playAudio() {
    const audio = document.getElementById('audio');
    const playBtn = document.getElementById('play-btn');

    audio.currentTime = 0;
    audio.play();

    playBtn.disabled = true;
    playBtn.textContent = '▶ Playing...';

    setTimeout(() => {
        audio.pause();
        playBtn.disabled = false;
        playBtn.textContent = '▶ Play';
    }, gameState.audioDuration * 1000);
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
