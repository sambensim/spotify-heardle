// Playlist selection page logic
document.addEventListener('DOMContentLoaded', async () => {
    const loading = document.getElementById('loading');
    const error = document.getElementById('error');
    const playlistsContainer = document.getElementById('playlists');

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
            card.onclick = () => selectPlaylist(playlist.id);

            const image = playlist.images && playlist.images.length > 0
                ? playlist.images[0].url
                : 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="150" height="150"><rect fill="%23ddd" width="150" height="150"/></svg>';

            card.innerHTML = `
                <img src="${image}" alt="${playlist.name}" class="playlist-image" onerror="this.src='data:image/svg+xml,<svg xmlns=\\'http://www.w3.org/2000/svg\\' width=\\'150\\' height=\\'150\\'><rect fill=\\'%23ddd\\' width=\\'150\\' height=\\'150\\'/></svg>'">
                <div class="playlist-name">${playlist.name}</div>
                <div class="playlist-tracks">${playlist.tracks.total} tracks</div>
            `;

            playlistsContainer.appendChild(card);
        });
    } catch (err) {
        loading.style.display = 'none';
        error.textContent = 'Failed to load playlists. Please try logging in again.';
        error.style.display = 'block';
        console.error('Error loading playlists:', err);
    }
});

function selectPlaylist(playlistId) {
    window.location.href = `/game.html?playlist=${playlistId}`;
}
