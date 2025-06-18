package domain

import "errors"

// Track representa uma faixa musical
type Track struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Duration    int    `json:"duration,omitempty"` // Duração em milissegundos
	Description string `json:"description,omitempty"`
	CoverURL    string `json:"coverUrl,omitempty"`
}

// Artwork represents the artwork information
type Artwork struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	URL        string `json:"url"`
	BgColor    string `json:"bgColor"`
	TextColor1 string `json:"textColor1"`
	TextColor2 string `json:"textColor2"`
	TextColor3 string `json:"textColor3"`
	TextColor4 string `json:"textColor4"`
}

// PlayParams represents the play parameters
type PlayParams struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

// Preview represents a preview URL
type Preview struct {
	URL string `json:"url"`
}

// EditorialNotes represents editorial notes
type EditorialNotes struct {
	Standard string `json:"standard,omitempty"`
	Short    string `json:"short,omitempty"`
}

// Song represents a song in the Apple Music catalog
type Song struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Href       string         `json:"href"`
	Attributes SongAttributes `json:"attributes"`
}

// SongAttributes represents the attributes of a song
type SongAttributes struct {
	AlbumName            string     `json:"albumName"`
	GenreNames           []string   `json:"genreNames"`
	TrackNumber          int        `json:"trackNumber"`
	ReleaseDate          string     `json:"releaseDate"`
	DurationInMillis     int        `json:"durationInMillis"`
	ISRC                 string     `json:"isrc"`
	Artwork              Artwork    `json:"artwork"`
	URL                  string     `json:"url"`
	PlayParams           PlayParams `json:"playParams"`
	DiscNumber           int        `json:"discNumber"`
	IsAppleDigitalMaster bool       `json:"isAppleDigitalMaster"`
	HasLyrics            bool       `json:"hasLyrics"`
	Name                 string     `json:"name"`
	Previews             []Preview  `json:"previews"`
	ArtistName           string     `json:"artistName"`
	ComposerName         string     `json:"composerName,omitempty"`
}

// Album represents an album in the Apple Music catalog
type Album struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Href       string          `json:"href"`
	Attributes AlbumAttributes `json:"attributes"`
}

// AlbumAttributes represents the attributes of an album
type AlbumAttributes struct {
	Copyright           string         `json:"copyright"`
	GenreNames          []string       `json:"genreNames"`
	ReleaseDate         string         `json:"releaseDate"`
	IsMasteredForItunes bool           `json:"isMasteredForItunes"`
	UPC                 string         `json:"upc"`
	Artwork             Artwork        `json:"artwork"`
	URL                 string         `json:"url"`
	PlayParams          PlayParams     `json:"playParams"`
	RecordLabel         string         `json:"recordLabel"`
	TrackCount          int            `json:"trackCount"`
	IsCompilation       bool           `json:"isCompilation"`
	IsSingle            bool           `json:"isSingle"`
	Name                string         `json:"name"`
	ArtistName          string         `json:"artistName"`
	EditorialNotes      EditorialNotes `json:"editorialNotes,omitempty"`
	IsComplete          bool           `json:"isComplete"`
	ContentRating       string         `json:"contentRating,omitempty"`
}

// Artist represents an artist in the Apple Music catalog
type Artist struct {
	ID            string              `json:"id"`
	Type          string              `json:"type"`
	Href          string              `json:"href"`
	Attributes    ArtistAttributes    `json:"attributes"`
	Relationships ArtistRelationships `json:"relationships,omitempty"`
}

// ArtistAttributes represents the attributes of an artist
type ArtistAttributes struct {
	Name       string   `json:"name"`
	GenreNames []string `json:"genreNames"`
	Artwork    Artwork  `json:"artwork"`
	URL        string   `json:"url"`
}

// ArtistRelationships represents the relationships of an artist
type ArtistRelationships struct {
	Albums struct {
		Href string `json:"href"`
		Data []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Href string `json:"href"`
		} `json:"data"`
	} `json:"albums"`
}

// SearchResults represents the search results in Apple Music API format
type SearchResults struct {
	Results struct {
		Artists *struct {
			Href string   `json:"href"`
			Data []Artist `json:"data"`
		} `json:"artists,omitempty"`
		Songs *struct {
			Href string `json:"href"`
			Next string `json:"next,omitempty"`
			Data []Song `json:"data"`
		} `json:"songs,omitempty"`
		Albums *struct {
			Href string  `json:"href"`
			Next string  `json:"next,omitempty"`
			Data []Album `json:"data"`
		} `json:"albums,omitempty"`
	} `json:"results"`
	Meta struct {
		Results struct {
			Order    []string `json:"order"`
			RawOrder []string `json:"rawOrder"`
		} `json:"results"`
	} `json:"meta"`
}

// Playlist representa uma playlist
type Playlist struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	CoverURL    string  `json:"coverUrl,omitempty"`
	Tracks      []Track `json:"tracks,omitempty"`
}

// MaxTracksPerPlaylist define o número máximo de faixas por playlist
const MaxTracksPerPlaylist = 100

// AddTrack adiciona uma faixa à playlist, respeitando as regras de negócio
func (p *Playlist) AddTrack(track Track) error {
	if len(p.Tracks) >= MaxTracksPerPlaylist {
		return errors.New("playlist has reached the maximum number of tracks")
	}

	p.Tracks = append(p.Tracks, track)
	return nil
}

// RemoveTrack remove uma faixa da playlist pelo seu ID
func (p *Playlist) RemoveTrack(trackID string) bool {
	initialLength := len(p.Tracks)
	var updatedTracks []Track
	for _, track := range p.Tracks {
		if track.ID != trackID {
			updatedTracks = append(updatedTracks, track)
		}
	}
	p.Tracks = updatedTracks
	return len(p.Tracks) < initialLength // Retorna true se uma faixa foi removida
}

// Station representa uma estação de rádio
type Station struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CoverURL    string `json:"coverUrl,omitempty"`
}

// Activity representa uma atividade musical
type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CoverURL    string `json:"coverUrl,omitempty"`
}

// Curator representa um curador de conteúdo
type Curator struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
}
