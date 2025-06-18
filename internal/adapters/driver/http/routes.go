package http

import (
	"applemusic-api-simulator/internal/core/ports/driving"
	"net/http"

	"github.com/gorilla/mux"
)

// Router configura as rotas da aplicação
func Router(searchHandler *SearchHandler) http.Handler {
	mux := http.NewServeMux()

	// Rota de busca
	mux.HandleFunc("/v1/catalog/us/search", searchHandler.Search)

	return mux
}

// loggingMiddleware registra informações sobre as requisições
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implementar logging adequado
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware configura os headers CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRoutes(router *mux.Router, musicService driving.MusicService) {
	searchHandler := NewSearchHandler(musicService)
	router.HandleFunc("/v1/catalog/us/search", searchHandler.Search).Methods("GET")
}
