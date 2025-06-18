package driving

import (
	"applemusic-api-simulator/internal/core/domain"
)

// SearchResultType representa os tipos de resultados de busca suportados
type SearchResultType string

const (
	ArtistsType SearchResultType = "artists"
	SongsType   SearchResultType = "songs"
	AlbumsType  SearchResultType = "albums"
)

// SearchParameters representa os parâmetros de uma busca
type SearchParameters struct {
	Term   string
	Limit  int
	Offset int
	Types  []SearchResultType
}

type SearchResults struct {
	Artists []domain.Artist `json:"artists"`
	Songs   []domain.Song   `json:"songs"`
	Albums  []domain.Album  `json:"albums"`
}

// MusicService define a interface para o serviço de música
type MusicService interface {
	// Search realiza uma busca por músicas, álbuns e artistas
	Search(params SearchParameters) (*SearchResults, error)
}
