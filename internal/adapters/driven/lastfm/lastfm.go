package lastfm

import (
	"applemusic-api-simulator/internal/core/domain"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const lastfmBaseURL = "https://ws.audioscrobbler.com/2.0/"

// LastFMAdapter implementa a interface MusicProvider
// Usa a API do Last.fm para buscar músicas, álbuns e artistas
// https://www.last.fm/api/show/track.search
// https://www.last.fm/api/show/album.search
// https://www.last.fm/api/show/artist.search
type LastFMAdapter struct {
	apiKey    string
	apiSecret string
	client    *http.Client
}

func NewLastFMAdapter() (*LastFMAdapter, error) {
	apiKey := os.Getenv("LASTFM_API_KEY")
	apiSecret := os.Getenv("LASTFM_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("LASTFM_API_KEY and LASTFM_API_SECRET environment variables are required")
	}

	return &LastFMAdapter{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    &http.Client{},
	}, nil
}

func (a *LastFMAdapter) SearchSongs(query string, limit, offset int) ([]domain.Song, error) {
	// Construir URL da API
	url := fmt.Sprintf("%s?method=track.search&track=%s&api_key=%s&format=json&limit=%d&page=%d",
		lastfmBaseURL,
		url.QueryEscape(query),
		a.apiKey,
		limit*2, // Buscar mais resultados para ter mais opções para filtrar
		offset/limit+1)

	// Fazer requisição
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Last.fm: %v", err)
	}
	defer resp.Body.Close()

	// Ler corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Decodificar resposta
	var result struct {
		Results struct {
			TrackMatches struct {
				Track []struct {
					Name   string `json:"name"`
					Artist string `json:"artist"`
				} `json:"track"`
			} `json:"trackmatches"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Converter para domínio
	var songs []domain.Song
	query = strings.ToLower(query)
	queryTerms := strings.Fields(query)

	// Primeiro, tentar encontrar músicas que correspondam exatamente ao termo de busca
	// e sejam da artista principal
	var exactMatches []domain.Song
	var otherMatches []domain.Song

	for _, track := range result.Results.TrackMatches.Track {
		trackName := strings.ToLower(track.Name)
		artistName := strings.ToLower(track.Artist)

		// Verificar se é uma música da artista principal
		isMainArtist := false
		for _, term := range queryTerms {
			if strings.Contains(artistName, term) && !strings.Contains(artistName, "(") && !strings.Contains(artistName, ")") {
				isMainArtist = true
				break
			}
		}

		// Verificar se o título corresponde exatamente
		exactTitleMatch := false
		for _, term := range queryTerms {
			if strings.Contains(trackName, term) {
				exactTitleMatch = true
				break
			}
		}

		songID := fmt.Sprintf("%s-%s", strings.ToLower(track.Artist), strings.ToLower(track.Name))
		songID = strings.ReplaceAll(songID, " ", "-")
		songID = strings.ReplaceAll(songID, "/", "-")
		songID = strings.ReplaceAll(songID, "\\", "-")

		song := domain.Song{
			ID:   songID,
			Type: "songs",
			Attributes: domain.SongAttributes{
				Name:       track.Name,
				ArtistName: track.Artist,
				GenreNames: []string{"Pop"},
			},
		}

		if isMainArtist && exactTitleMatch {
			exactMatches = append(exactMatches, song)
		} else {
			otherMatches = append(otherMatches, song)
		}
	}

	// Combinar resultados, priorizando matches exatos
	songs = append(songs, exactMatches...)
	songs = append(songs, otherMatches...)

	// Limitar ao número solicitado
	if len(songs) > limit {
		songs = songs[:limit]
	}

	return songs, nil
}

func (a *LastFMAdapter) SearchAlbums(term string, limit, offset int) ([]domain.Album, error) {
	// Construir URL da API
	baseURL := "http://ws.audioscrobbler.com/2.0/"
	params := url.Values{}
	params.Set("method", "album.search")
	params.Set("album", term)
	params.Set("api_key", a.apiKey)
	params.Set("format", "json")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("page", strconv.Itoa(offset/limit+1))

	// Fazer requisição
	resp, err := a.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("error making request to Last.fm: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Decodificar resposta
	var result struct {
		Results struct {
			AlbumMatches struct {
				Album []struct {
					Name      string `json:"name"`
					Artist    string `json:"artist"`
					URL       string `json:"url"`
					Listeners string `json:"listeners"`
					Image     []struct {
						Size string `json:"size"`
						URL  string `json:"#text"`
					} `json:"image"`
				} `json:"album"`
			} `json:"albummatches"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Converter para domínio
	var albums []domain.Album
	for _, album := range result.Results.AlbumMatches.Album {
		// Encontrar imagem de capa
		var artworkURL string
		for _, img := range album.Image {
			if img.Size == "large" {
				artworkURL = img.URL
				break
			}
		}

		// Gerar ID único baseado no nome e artista
		id := fmt.Sprintf("%s-%s", strings.ToLower(album.Artist), strings.ToLower(album.Name))
		id = strings.ReplaceAll(id, " ", "-")
		id = strings.ReplaceAll(id, "/", "-")
		id = strings.ReplaceAll(id, "\\", "-")

		album := domain.Album{
			ID:   id,
			Type: "albums",
			Href: fmt.Sprintf("/v1/catalog/us/albums/%s", id),
			Attributes: domain.AlbumAttributes{
				Name:       album.Name,
				ArtistName: album.Artist,
				Artwork: domain.Artwork{
					URL: artworkURL,
				},
				PlayParams: domain.PlayParams{
					ID:   id,
					Kind: "album",
				},
				GenreNames: []string{"Pop"},
				IsComplete: true,
			},
		}
		albums = append(albums, album)

		// Limitar número de resultados
		if len(albums) >= 5 {
			break
		}
	}

	return albums, nil
}

func (a *LastFMAdapter) SearchArtists(term string, limit, offset int) ([]domain.Artist, error) {
	// Construir URL da API
	baseURL := "http://ws.audioscrobbler.com/2.0/"
	params := url.Values{}
	params.Set("method", "artist.search")
	params.Set("artist", term)
	params.Set("api_key", a.apiKey)
	params.Set("format", "json")
	params.Set("limit", strconv.Itoa(limit))
	params.Set("page", strconv.Itoa(offset/limit+1))

	// Fazer requisição
	resp, err := a.client.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, fmt.Errorf("error making request to Last.fm: %w", err)
	}
	defer resp.Body.Close()

	// Ler resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Decodificar resposta
	var result struct {
		Results struct {
			ArtistMatches struct {
				Artist []struct {
					Name      string `json:"name"`
					URL       string `json:"url"`
					Listeners string `json:"listeners"`
					Image     []struct {
						Size string `json:"size"`
						URL  string `json:"#text"`
					} `json:"image"`
				} `json:"artist"`
			} `json:"artistmatches"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Converter para domínio
	var artists []domain.Artist
	for _, artist := range result.Results.ArtistMatches.Artist {
		// Encontrar imagem do artista
		var artworkURL string
		for _, img := range artist.Image {
			if img.Size == "large" {
				artworkURL = img.URL
				break
			}
		}

		// Gerar ID único baseado no nome
		id := strings.ToLower(artist.Name)
		id = strings.ReplaceAll(id, " ", "-")
		id = strings.ReplaceAll(id, "/", "-")
		id = strings.ReplaceAll(id, "\\", "-")

		artist := domain.Artist{
			ID:   id,
			Type: "artists",
			Href: fmt.Sprintf("/v1/catalog/us/artists/%s", id),
			Attributes: domain.ArtistAttributes{
				Name:       artist.Name,
				GenreNames: []string{"Pop"},
				Artwork: domain.Artwork{
					URL: artworkURL,
				},
			},
		}
		artists = append(artists, artist)

		// Limitar número de resultados
		if len(artists) >= 5 {
			break
		}
	}

	return artists, nil
}

// Implementações vazias para os outros tipos de busca
func (a *LastFMAdapter) searchPlaylists(term string, limit, offset int) ([]domain.Playlist, error) {
	// Last.fm não tem API para playlists, retornando lista vazia
	return []domain.Playlist{}, nil
}

func (a *LastFMAdapter) searchStations(term string, limit, offset int) ([]domain.Station, error) {
	// Last.fm não tem API para estações, retornando lista vazia
	return []domain.Station{}, nil
}

func (a *LastFMAdapter) searchActivities(term string, limit, offset int) ([]domain.Activity, error) {
	// Last.fm não tem API para atividades, retornando lista vazia
	return []domain.Activity{}, nil
}

func (a *LastFMAdapter) searchCurators(term string, limit, offset int) ([]domain.Curator, error) {
	// Last.fm não tem API para curadores, retornando lista vazia
	return []domain.Curator{}, nil
}
