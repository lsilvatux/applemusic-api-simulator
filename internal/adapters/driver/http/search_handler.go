package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"applemusic-api-simulator/internal/core/domain"
	"applemusic-api-simulator/internal/core/ports/driving"
)

// SearchHandler lida com as requisições de busca
type SearchHandler struct {
	musicService driving.MusicService
}

// NewSearchHandler cria uma nova instância do handler de busca
func NewSearchHandler(musicService driving.MusicService) *SearchHandler {
	return &SearchHandler{
		musicService: musicService,
	}
}

// Search processa a requisição de busca
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Obter parâmetros da query string
	term := r.URL.Query().Get("term")
	if term == "" {
		http.Error(w, "term parameter is required", http.StatusBadRequest)
		return
	}

	// Obter e validar limit
	limitStr := r.URL.Query().Get("limit")
	limit := 5 // Default limit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit parameter", http.StatusBadRequest)
			return
		}
		if limit < 1 {
			limit = 5
		} else if limit > 25 {
			limit = 25
		}
	}

	// Obter e validar offset
	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset parameter", http.StatusBadRequest)
			return
		}
		if offset < 0 {
			offset = 0
		}
	}

	// Obter e validar types
	typesStr := r.URL.Query().Get("types")
	var types []driving.SearchResultType
	if typesStr != "" {
		typeStrs := strings.Split(typesStr, ",")
		for _, t := range typeStrs {
			switch strings.TrimSpace(t) {
			case "artists":
				types = append(types, driving.ArtistsType)
			case "songs":
				types = append(types, driving.SongsType)
			case "albums":
				types = append(types, driving.AlbumsType)
			}
		}
	} else {
		// Default to all types if none specified
		types = []driving.SearchResultType{
			driving.ArtistsType,
			driving.SongsType,
			driving.AlbumsType,
		}
	}

	// Construir parâmetros de busca
	params := driving.SearchParameters{
		Term:   term,
		Limit:  limit,
		Offset: offset,
		Types:  types,
	}

	// Realizar busca
	results, err := h.musicService.Search(params)
	if err != nil {
		http.Error(w, fmt.Sprintf("error performing search: %v", err), http.StatusInternalServerError)
		return
	}

	// Construir resposta no formato da Apple Music API
	response := struct {
		Results map[string]interface{} `json:"results"`
		Meta    struct {
			Results struct {
				Order []string `json:"order"`
			} `json:"results"`
		} `json:"meta"`
	}{
		Results: make(map[string]interface{}),
		Meta: struct {
			Results struct {
				Order []string `json:"order"`
			} `json:"results"`
		}{
			Results: struct {
				Order []string `json:"order"`
			}{
				Order: make([]string, 0, len(types)),
			},
		},
	}

	// Adicionar apenas os tipos solicitados à resposta
	for _, t := range types {
		switch t {
		case driving.ArtistsType:
			response.Results["artists"] = struct {
				Href string          `json:"href"`
				Next string          `json:"next"`
				Data []domain.Artist `json:"data"`
			}{
				Href: fmt.Sprintf("/v1/catalog/us/search?term=%s&types=artists&limit=%d&offset=%d", term, limit, offset),
				Next: fmt.Sprintf("/v1/catalog/us/search?term=%s&types=artists&limit=%d&offset=%d", term, limit, offset+limit),
				Data: results.Artists,
			}
			response.Meta.Results.Order = append(response.Meta.Results.Order, "artists")

		case driving.SongsType:
			response.Results["songs"] = struct {
				Href string        `json:"href"`
				Next string        `json:"next"`
				Data []domain.Song `json:"data"`
			}{
				Href: fmt.Sprintf("/v1/catalog/us/search?term=%s&types=songs&limit=%d&offset=%d", term, limit, offset),
				Next: fmt.Sprintf("/v1/catalog/us/search?term=%s&types=songs&limit=%d&offset=%d", term, limit, offset+limit),
				Data: results.Songs,
			}
			response.Meta.Results.Order = append(response.Meta.Results.Order, "songs")

		case driving.AlbumsType:
			response.Results["albums"] = struct {
				Href string         `json:"href"`
				Next string         `json:"next"`
				Data []domain.Album `json:"data"`
			}{
				Href: fmt.Sprintf("/v1/catalog/us/search?term=%s&types=albums&limit=%d&offset=%d", term, limit, offset),
				Next: fmt.Sprintf("/v1/catalog/us/search?term=%s&types=albums&limit=%d&offset=%d", term, limit, offset+limit),
				Data: results.Albums,
			}
			response.Meta.Results.Order = append(response.Meta.Results.Order, "albums")
		}
	}

	// Enviar resposta
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}
