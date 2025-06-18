package services

import (
	"fmt"

	"applemusic-api-simulator/internal/core/ports/driven"
	"applemusic-api-simulator/internal/core/ports/driving"
)

type MusicService struct {
	musicProvider driven.MusicProvider
}

func NewMusicService(musicProvider driven.MusicProvider) driving.MusicService {
	return &MusicService{
		musicProvider: musicProvider,
	}
}

func (s *MusicService) Search(params driving.SearchParameters) (*driving.SearchResults, error) {
	results := &driving.SearchResults{}

	// Para cada tipo solicitado, realizar a busca específica
	for _, searchType := range params.Types {
		switch searchType {
		case driving.SongsType:
			songs, err := s.musicProvider.SearchSongs(params.Term, params.Limit, params.Offset)
			if err != nil {
				return nil, fmt.Errorf("error searching songs: %w", err)
			}
			// Garantir que não exceda o limite
			if len(songs) > params.Limit {
				songs = songs[:params.Limit]
			}
			results.Songs = songs

		case driving.AlbumsType:
			albums, err := s.musicProvider.SearchAlbums(params.Term, params.Limit, params.Offset)
			if err != nil {
				return nil, fmt.Errorf("error searching albums: %w", err)
			}
			// Garantir que não exceda o limite
			if len(albums) > params.Limit {
				albums = albums[:params.Limit]
			}
			results.Albums = albums

		case driving.ArtistsType:
			artists, err := s.musicProvider.SearchArtists(params.Term, params.Limit, params.Offset)
			if err != nil {
				return nil, fmt.Errorf("error searching artists: %w", err)
			}
			// Garantir que não exceda o limite
			if len(artists) > params.Limit {
				artists = artists[:params.Limit]
			}
			results.Artists = artists
		}
	}

	return results, nil
}
