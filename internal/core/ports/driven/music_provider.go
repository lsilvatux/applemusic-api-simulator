package driven

import (
	"applemusic-api-simulator/internal/core/domain"
)

// ProviderSearchResults contém os resultados de uma busca diretamente do provedor de dados.
// Pode ser muito similar ou idêntico ao 'driving.SearchResults' se o provedor
// já retorna os dados de uma forma que podemos mapear diretamente para o nosso domínio.
// Ou poderia ser uma estrutura mais genérica se precisássemos de mais processamento.
// Por simplicidade inicial, vamos mantê-lo similar.
type ProviderSearchResults struct {
	Tracks  []domain.Track
	Albums  []domain.Album
	Artists []domain.Artist
}

// MusicProvider define a interface para provedores de música
type MusicProvider interface {
	// SearchSongs busca músicas com base no termo de busca
	SearchSongs(term string, limit, offset int) ([]domain.Song, error)

	// SearchAlbums busca álbuns com base no termo de busca
	SearchAlbums(term string, limit, offset int) ([]domain.Album, error)

	// SearchArtists busca artistas com base no termo de busca
	SearchArtists(term string, limit, offset int) ([]domain.Artist, error)
}
