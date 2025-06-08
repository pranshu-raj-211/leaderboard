# Real-time Game Leaderboard

A high-performance real-time game leaderboard system built with Go, Redis, and Server-Sent Events (SSE).

## Features

- Real-time leaderboard updates using SSE
- Redis-backed scoring system
- REST API endpoints for game submissions and statistics
- Prometheus metrics integration
- Structured logging with Zap
- Configuration management with YAML

## Prerequisites

- Go 1.22 or higher
- Redis (version 8)
- Docker (and compose)

## Quick Start

1. Clone the repository:

```bash
git clone 
cd leaderboard
```

2. Install dependencies (for dev):

```bash
go mod download
```

3. Configure the application in `config.yaml`:


4. Run the application:

```bash
docker compose up --build
```

## API Endpoints

### Submit Game Results

```http
POST /submit-game
Content-Type: application/json

{
  "game_id": "game123",
  "server_id": "server1",
  "player1_id": "player1",
  "player2_id": "player2",
  "result": 1
}
```

### Get Leaderboard

```http
GET /leaderboard
```

### Get Player Stats

```http
GET /player/:id/stats
```

### Stream Real-time Updates

```http
GET /stream-leaderboard
```

### Prometheus Metrics

```http
GET /metrics
```

## Monitoring

The application exports Prometheus metrics at `/metrics` including:

- Game submission counts
- HTTP request durations
- Redis operation latencies
- Error counts

## Development

### Project Structure

```
leaderboard/
├── src/
│   ├── backend/     # API handlers
│   ├── config/      # Configuration and logging
│   ├── metrics/     # Prometheus metrics
│   ├── models/      # Data models
│   ├── redisclient/ # Redis operations
│   └── main.go      # Application entry point
├── config.yaml      # Configuration file
├── go.mod          # Go module file
└── README.md       # This file
```

## TODO

- [ ] Add rate limiting for API endpoints
- [ ] Add authentication/authorization (necessary for server)
- [ ] Create comprehensive API documentation
- [ ] Add tests
- [ ] Implement data persistence backup (Postgres, currently in the works)
- [ ] Add request validation middleware
- [ ] Create Grafana dashboards (in progress)
- [ ] Add player history tracking (and translate ids to usernames before sending)