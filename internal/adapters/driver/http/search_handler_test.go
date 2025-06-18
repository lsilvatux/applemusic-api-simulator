package http

import (
	"applemusic-api-simulator/internal/core/domain"
	"applemusic-api-simulator/internal/core/ports/driven"
	"applemusic-api-simulator/internal/core/ports/driving"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type AppleMusicResponse struct {
	Results struct {
		Songs struct {
			Data []AppleMusicTrack `json:"data"`
		} `json:"songs"`
		Albums struct {
			Data []AppleMusicAlbum `json:"data"`
		} `json:"albums"`
		Artists struct {
			Data []AppleMusicArtist `json:"data"`
		} `json:"artists"`
	} `json:"results"`
}

type AppleMusicTrack struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name             string `json:"name"`
		ArtistName       string `json:"artistName"`
		DurationInMillis int    `json:"durationInMillis"`
		Description      string `json:"description,omitempty"`
	} `json:"attributes"`
}

type AppleMusicAlbum struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name        string `json:"name"`
		ArtistName  string `json:"artistName"`
		Description string `json:"description,omitempty"`
	} `json:"attributes"`
}

type AppleMusicArtist struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		Genre       string `json:"genre,omitempty"`
	} `json:"attributes"`
}

// mockMusicProvider é um mock do MusicProvider para testes
type mockMusicProvider struct {
	searchResults *driven.ProviderSearchResults
	searchError   error
}

func (m *mockMusicProvider) Search(params driving.SearchParameters) (*driving.SearchResults, error) {
	return &driving.SearchResults{
		Artists: []domain.Artist{},
		Songs:   []domain.Song{},
		Albums:  []domain.Album{},
	}, m.searchError
}

func TestSearchHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		types          string
		mockResults    *driven.ProviderSearchResults
		mockError      error
		expectedStatus int
		expectedBody   *AppleMusicResponse
	}{
		{
			name:  "Successful search with all types",
			query: "test",
			types: "",
			mockResults: &driven.ProviderSearchResults{
				Tracks: []domain.Track{
					{
						ID:          "1",
						Title:       "Test Track",
						Artist:      "Test Artist",
						Duration:    180000,
						Description: "Test Description",
					},
				},
				Albums: []domain.Album{
					{
						ID:   "2",
						Type: "albums",
						Attributes: domain.AlbumAttributes{
							Name:       "Test Album",
							ArtistName: "Test Artist",
							EditorialNotes: domain.EditorialNotes{
								Standard: "Test Description",
							},
						},
					},
				},
				Artists: []domain.Artist{
					{
						ID:   "3",
						Type: "artists",
						Attributes: domain.ArtistAttributes{
							Name:       "Test Artist",
							GenreNames: []string{"Pop"},
						},
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: &AppleMusicResponse{
				Results: struct {
					Songs struct {
						Data []AppleMusicTrack `json:"data"`
					} `json:"songs"`
					Albums struct {
						Data []AppleMusicAlbum `json:"data"`
					} `json:"albums"`
					Artists struct {
						Data []AppleMusicArtist `json:"data"`
					} `json:"artists"`
				}{
					Songs: struct {
						Data []AppleMusicTrack `json:"data"`
					}{
						Data: []AppleMusicTrack{
							{
								ID:   "1",
								Type: "songs",
								Attributes: struct {
									Name             string `json:"name"`
									ArtistName       string `json:"artistName"`
									DurationInMillis int    `json:"durationInMillis"`
									Description      string `json:"description,omitempty"`
								}{
									Name:             "Test Track",
									ArtistName:       "Test Artist",
									DurationInMillis: 180000,
									Description:      "Test Description",
								},
							},
						},
					},
					Albums: struct {
						Data []AppleMusicAlbum `json:"data"`
					}{
						Data: []AppleMusicAlbum{
							{
								ID:   "2",
								Type: "albums",
								Attributes: struct {
									Name        string `json:"name"`
									ArtistName  string `json:"artistName"`
									Description string `json:"description,omitempty"`
								}{
									Name:        "Test Album",
									ArtistName:  "Test Artist",
									Description: "Test Description",
								},
							},
						},
					},
					Artists: struct {
						Data []AppleMusicArtist `json:"data"`
					}{
						Data: []AppleMusicArtist{
							{
								ID:   "3",
								Type: "artists",
								Attributes: struct {
									Name        string `json:"name"`
									Description string `json:"description,omitempty"`
									Genre       string `json:"genre,omitempty"`
								}{
									Name:        "Test Artist",
									Description: "Test Description",
									Genre:       "Pop",
								},
							},
						},
					},
				},
			},
		},
		{
			name:           "Missing term parameter",
			query:          "",
			types:          "",
			mockResults:    nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:           "Provider error",
			query:          "test",
			types:          "",
			mockResults:    nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar mock do provider
			mockProvider := &mockMusicProvider{
				searchResults: tt.mockResults,
				searchError:   tt.mockError,
			}

			// Criar handler
			handler := NewSearchHandler(mockProvider)

			// Criar request
			req := httptest.NewRequest("GET", "/v1/catalog/us/search", nil)
			q := req.URL.Query()
			if tt.query != "" {
				q.Add("term", tt.query)
			}
			if tt.types != "" {
				q.Add("types", tt.types)
			}
			req.URL.RawQuery = q.Encode()

			// Criar response recorder
			rr := httptest.NewRecorder()

			// Executar handler
			handler.Search(rr, req)

			// Verificar status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, rr.Code)
			}

			// Se esperamos um corpo de resposta
			if tt.expectedBody != nil {
				var response AppleMusicResponse
				if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
					t.Fatalf("error decoding response: %v", err)
				}

				// Verificar se a resposta corresponde ao esperado
				if !reflect.DeepEqual(&response, tt.expectedBody) {
					t.Errorf("expected response %+v; got %+v", tt.expectedBody, response)
				}
			}
		})
	}
}
