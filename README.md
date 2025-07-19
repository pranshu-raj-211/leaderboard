# Real-time Game Leaderboard

A high-performance real-time game leaderboard system built with Go, Redis, and Server-Sent Events (SSE).

## Features

- Real-time leaderboard updates using SSE
- Redis-backed scoring system
- REST API endpoints for game submissions and statistics
- Prometheus metrics integration
- Structured logging with Zap
- Configuration management with YAML
- Docker compose based
- Grafana based dashboard (configured through `yaml` and `json` - no setup needed)

## Prerequisites

- Docker (and Docker Compose)

For the dev environment it's good to have Go 1.24 installed.

## Quick Start

1. Clone the repository:

```bash
git clone 
cd leaderboard
```

2. Run the application:

```bash
docker compose up --build
```

Dev: Install dependencies

```bash
go mod download
```

(Optional) Configure the application in `config.yaml`:

To checkout the Grafana dashboard, run the app and go to `http://localhost:3000`, login with the default Grafan username and password (admin).

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

and many more. Check out the grafana dashboard at `http://localhost:3000` after running with docker compose. Username and password are grafana defaults (admin).

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
- [ ] Add player history tracking (and translate ids to usernames before sending)

## Ideas to check out

1. Builder images and go image vulnerabilities.
2. Securing docker containers
3. Profiling heap and CPU (pprof) in detail
