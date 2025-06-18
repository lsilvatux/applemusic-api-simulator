# Apple Music API Simulator

This is a Go-based simulator that mimics the Apple Music API search endpoint, built using hexagonal architecture principles. This project provides a REST API that accepts search queries and returns results in Apple Music's response format, using Last.fm as the underlying music data provider.

I created this project to learn about the hexagonal architecture and how to implement it in GoLang.

## Features

- **Apple Music API Compatibility**: Simulates the Apple Music `/v1/catalog/us/search` endpoint
- **Hexagonal Architecture**: Clean separation of concerns with adapters, ports, and domain logic
- **Last.fm Integration**: Uses Last.fm API as the music data provider
- **Search Support**: Search for tracks, albums, and artists
- **Query Parameters**: Supports `term`, `types`, `limit`, and `offset` parameters
- **Response Format**: Returns results in Apple Music's JSON response format

## Architecture

The project follows hexagonal architecture principles:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Handler (Driver)                    │
├─────────────────────────────────────────────────────────────┤
│                    Search Service                           │
├─────────────────────────────────────────────────────────────┤
│                    Domain Entities                          │
├─────────────────────────────────────────────────────────────┤
│                    Last.fm Adapter (Driven)                 │
└─────────────────────────────────────────────────────────────┘
```

- **Driver Adapters**: HTTP handlers that expose the API
- **Domain**: Core business logic and entities
- **Driven Adapters**: External service integrations (Last.fm)
- **Ports**: Interfaces defining contracts between layers

## Prerequisites

- Go 1.21 or higher
- Last.fm API credentials (API key and shared secret)

## Environment Variables

Create a `.env` file in the root directory with your Last.fm credentials:

```env
LASTFM_API_KEY=your_api_key_here
LASTFM_SHARED_SECRET=your_shared_secret_here
```

To get Last.fm API credentials:
1. Visit [Last.fm API](https://www.last.fm/api)
2. Create an account and register your application
3. Get your API key and shared secret

## Installation and Running

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd applemusic-api-simulator
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment variables** (create `.env` file as shown above)

4. **Run the server**:
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:8080`

## API Usage

### Search Endpoint

**Endpoint**: `GET /v1/catalog/us/search`

**Query Parameters**:
- `term` (required): Search term
- `types` (optional): Comma-separated list of types to search for (`songs`, `albums`, `artists`)
- `limit` (optional): Number of results per type (default: 5, max: 25)
- `offset` (optional): Number of results to skip (default: 0)

### Examples

#### 1. Basic Search
Search for "The Beatles" across all types:

```bash
curl "http://localhost:8080/v1/catalog/us/search?term=The%20Beatles"
```

#### 2. Search Specific Types
Search only for songs and albums:

```bash
curl "http://localhost:8080/v1/catalog/us/search?term=Queen&types=songs,albums"
```

#### 3. Search with Limit
Limit results to 10 per type:

```bash
curl "http://localhost:8080/v1/catalog/us/search?term=Michael%20Jackson&limit=10"
```

#### 4. Search with Offset
Skip first 5 results:

```bash
curl "http://localhost:8080/v1/catalog/us/search?term=Rock&offset=5"
```

#### 5. Combined Parameters
Search for artists only with custom limit and offset:

```bash
curl "http://localhost:8080/v1/catalog/us/search?term=Pop&types=artists&limit=15&offset=10"
```

### Response Format

The API returns responses in Apple Music's format:

```json
{
  "results": {
    "songs": {
      "data": [
        {
          "id": "track_id",
          "type": "songs",
          "attributes": {
            "name": "Track Name",
            "artistName": "Artist Name",
            "durationInMillis": 180000,
            "description": "Track description"
          }
        }
      ]
    },
    "albums": {
      "data": [
        {
          "id": "album_id",
          "type": "albums",
          "attributes": {
            "name": "Album Name",
            "artistName": "Artist Name",
            "description": "Album description"
          }
        }
      ]
    },
    "artists": {
      "data": [
        {
          "id": "artist_id",
          "type": "artists",
          "attributes": {
            "name": "Artist Name",
            "description": "Artist description",
            "genre": "Primary genre"
          }
        }
      ]
    }
  }
}
```

## Error Responses

### 400 Bad Request
When required parameters are missing:
```json
{
  "error": "term parameter is required"
}
```

### 500 Internal Server Error
When the music provider encounters an error:
```json
{
  "error": "failed to search music"
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Project Structure

```
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── adapters/
│   │   ├── driver/
│   │   │   └── http/
│   │   │       ├── search_handler.go
│   │   │       └── search_handler_test.go
│   │   └── driven/
│   │       └── lastfm/
│   │           ├── lastfm.go
│   │           └── lastfm_test.go
│   ├── core/
│   │   ├── domain/
│   │   │   └── entities.go     # Domain entities
│   │   └── ports/
│   │       ├── driven/
│   │       │   └── music_provider.go
│   │       └── driving/
│   │           └── search_service.go
│   └── services/
│       └── search_service.go   # Business logic
├── go.mod
├── go.sum
└── README.md
```

## Acknowledgments

- [Last.fm API](https://www.last.fm/api) for providing music data
- Apple Music API for the response format specification 