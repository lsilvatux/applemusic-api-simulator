package main

import (
	"applemusic-api-simulator/internal/adapters/driven/lastfm"
	httpadapter "applemusic-api-simulator/internal/adapters/driver/http"
	"applemusic-api-simulator/internal/core/services"
	"log"
	"net/http"
)

func main() {
	// Inicializar o adaptador Last.fm
	lastfmAdapter, err := lastfm.NewLastFMAdapter()
	if err != nil {
		log.Fatalf("Error creating Last.fm adapter: %v", err)
	}

	// Inicializar o serviço de música
	musicService := services.NewMusicService(lastfmAdapter)

	// Inicializar o handler de busca
	searchHandler := httpadapter.NewSearchHandler(musicService)

	// Configurar as rotas
	router := httpadapter.Router(searchHandler)

	// Iniciar o servidor
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
