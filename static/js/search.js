// Search functionality with debouncing
let searchTimeout;
let currentSearchResults = [];

function initSearch() {
    const searchInput = document.getElementById('search-input');
    const searchResults = document.getElementById('search-results');

    searchInput.addEventListener('input', (e) => {
        const query = e.target.value.trim();

        if (query.length < 2) {
            searchResults.innerHTML = '';
            searchResults.style.display = 'none';
            return;
        }

        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => performSearch(query), 300);
    });

    searchInput.addEventListener('focus', () => {
        if (currentSearchResults.length > 0) {
            searchResults.style.display = 'block';
        }
    });

    document.addEventListener('click', (e) => {
        if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
            searchResults.style.display = 'none';
        }
    });
}

async function performSearch(query) {
    const searchResults = document.getElementById('search-results');
    
    try {
        const tracks = await searchTracks(query);
        currentSearchResults = tracks;
        
        searchResults.innerHTML = '';
        
        if (tracks.length === 0) {
            searchResults.innerHTML = '<div class="search-result-item">No results found</div>';
            searchResults.style.display = 'block';
            return;
        }

        tracks.forEach(track => {
            const item = document.createElement('div');
            item.className = 'search-result-item';
            item.onclick = () => selectTrack(track);

            const artists = track.artists.join(', ');
            
            item.innerHTML = `
                <div class="result-name">${track.name}</div>
                <div class="result-artist">${artists}</div>
            `;

            searchResults.appendChild(item);
        });

        searchResults.style.display = 'block';
    } catch (error) {
        console.error('Search failed:', error);
        searchResults.innerHTML = '<div class="search-result-item">Search failed</div>';
        searchResults.style.display = 'block';
    }
}

function selectTrack(track) {
    const searchResults = document.getElementById('search-results');
    searchResults.style.display = 'none';
    
    handleGuess(track.id, track.name);
}
