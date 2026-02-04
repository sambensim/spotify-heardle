// Spotify Heardle Clone - A music guessing game using Spotify playlists
package main

import (
	"log"
	"net/http"
	"spotify-heardle/config"
	"spotify-heardle/handlers"
	"spotify-heardle/storage"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	store := storage.NewMemoryStore()

	authHandler := handlers.NewAuthHandler(cfg, store)
	playlistHandler := handlers.NewPlaylistHandler(authHandler)
	searchHandler := handlers.NewSearchHandler(authHandler)
	gameHandler := handlers.NewGameHandler(authHandler, store)

	mux := http.NewServeMux()

	mux.HandleFunc("/login", authHandler.HandleLogin)
	mux.HandleFunc("/callback", authHandler.HandleCallback)
	mux.HandleFunc("/api/logout", authHandler.HandleLogout)
	mux.HandleFunc("/api/playlists", playlistHandler.HandleGetPlaylists)
	mux.HandleFunc("/api/search", searchHandler.HandleSearch)
	mux.HandleFunc("/api/game/start", gameHandler.HandleStartGame)
	mux.HandleFunc("/api/game/guess", gameHandler.HandleSubmitGuess)
	mux.HandleFunc("/api/game/skip", gameHandler.HandleSkip)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", fs)

	corsHandler := corsMiddleware(mux)
	loggedHandler := loggingMiddleware(corsHandler)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, loggedHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
