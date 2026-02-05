// Playlist selection page logic
const selectedPlaylists = new Set();

document.addEventListener('DOMContentLoaded', async () => {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const playlistsContainer = document.getElementById('playlists');
    const startButtonContainer = document.getElementById('start-button-container');

    try {
        const playlists = await getPlaylists();
        
        loading.style.display = 'none';

        if (playlists.length === 0) {
            error.textContent = 'No playlists found. Please create some playlists in Spotify.';
            error.style.display = 'block';
            return;
        }

        playlists.forEach(playlist => {
            const card = document.createElement('div');
            card.className = 'playlist-card';
            card.setAttribute('data-playlist-id', playlist.id);

            const image = playlist.images && playlist.images.length > 0
                ? playlist.images[0].url
                : 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="150" height="150"><rect fill="%23ddd" width="150" height="150"/></svg>';

            card.innerHTML = `
                <div class="playlist-checkbox-container">
                    <input type="checkbox" class="playlist-checkbox" id="playlist-${playlist.id}" data-playlist-id="${playlist.id}">
                </div>
                <img src="${image}" alt="${playlist.name}" class="playlist-image" onerror="this.src='data:image/svg+xml,<svg xmlns=\\'http://www.w3.org/2000/svg\\' width=\\'150\\' height=\\'150\\'><rect fill=\\'%23ddd\\' width=\\'150\\' height=\\'150\\'/></svg>'">
                <div class="playlist-name">${playlist.name}</div>
                <div class="playlist-tracks">${playlist.tracks.total} tracks</div>
            `;

            // Toggle selection when clicking the card
            card.addEventListener('click', (e) => {
                // Don't toggle if clicking directly on the checkbox
                if (e.target.classList.contains('playlist-checkbox')) {
                    return;
                }
                const checkbox = card.querySelector('.playlist-checkbox');
                checkbox.checked = !checkbox.checked;
                togglePlaylistSelection(playlist.id, checkbox.checked);
            });

            // Handle checkbox change
            const checkbox = card.querySelector('.playlist-checkbox');
            checkbox.addEventListener('change', (e) => {
                e.stopPropagation();
                togglePlaylistSelection(playlist.id, e.target.checked);
            });

            playlistsContainer.appendChild(card);
        });

        // Show start button container
        startButtonContainer.style.display = 'block';
    } catch (err) {
        loading.style.display = 'none';
        error.textContent = 'Failed to load playlists. Please try logging in again.';
        error.style.display = 'block';
        console.error('Error loading playlists:', err);
    }
});

function togglePlaylistSelection(playlistId, isSelected) {
    const card = document.querySelector(`[data-playlist-id="${playlistId}"]`);
    
    if (isSelected) {
        selectedPlaylists.add(playlistId);
        card.classList.add('selected');
    } else {
        selectedPlaylists.delete(playlistId);
        card.classList.remove('selected');
    }
    
    updateStartButton();
}

function updateStartButton() {
    const startBtn = document.getElementById('start-game-btn');
    const count = selectedPlaylists.size;
    
    if (count === 0) {
        startBtn.textContent = 'Start Game';
        startBtn.disabled = true;
    } else {
        startBtn.textContent = `Start Game with ${count} Playlist${count > 1 ? 's' : ''}`;
        startBtn.disabled = false;
    }
}

function startGameWithSelected() {
    if (selectedPlaylists.size === 0) {
        alert('Please select at least one playlist');
        return;
    }
    
    const playlistIds = Array.from(selectedPlaylists);
    window.location.href = `/game.html?playlists=${encodeURIComponent(JSON.stringify(playlistIds))}`;
}
