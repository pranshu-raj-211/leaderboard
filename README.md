# Real-time Game Leaderboard

Blogs on this project
- [28k+ conns, zero messages](https://blog.pranshu-raj.me/posts/implementing-correct-fanout)
- [Optimizing docker image builds](https://blog.pranshu-raj.me/posts/optimizing-docker-builds)
- (Upcoming) [Backpressure](https://blog.pranshu-raj.me/posts/understanding-backpressure/)


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

To checkout the Grafana dashboard, run the app and go to `http://localhost:3000`, login with the default Grafana username and password (admin).

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
|── tests/
|   |──load/         # Load testing scripts
├── config.yaml      # Configuration file
├── go.mod           # Go module file
└── README.md        # This file
```

## Benchmarking
To load test the application, first spin up the app using docker compose.

`docker compose up --build`

Then in a separate shell, after the compose app is up and running:

```bash
# POST requests done through HTTP
go run tests/load/submitter/submit.go

# and in another shell, for streaming responses - SSE clients
go run tests/load/streamer/stream.go
```

By default these will run with the default number of submitters as 10, SSE clients as 20000. On a Fedora Linux machine with 8GB RAM, this runs comfortably, with memory usage maxing out at 6.8GB.

[Post with Grafana dashboard - 28231 active SSE conns](https://x.com/seigino99707047/status/1955225344744583581)

[Post with game submissions dashboard](https://x.com/seigino99707047/status/1955246776035721718)

Combined streamer and submitter scripts have not been tested beyond this limit, however the streamer was tested in standalone mode, which went to 28,232 connections (default max outbound ports in Linux).

The streaming script is based off of [Eran Yanay's repo for his Gophercon talk](https://github.com/eranyanay/1m-go-websockets).

## TODO

- [ ] Add authentication (necessary for server)
- [ ] Create comprehensive API documentation
- [ ] Tests
- [ ] Implement data persistence backup and aggregation (Postgres, currently in the works)
- [ ] Add request validation middleware
- [ ] Add player history tracking (and translate ids to usernames before sending)
- [ ] Improve metrics and dashboard configs (quantiles)
- [ ] Redis pipelining and connection pooling (improve other access patterns)

## Ideas to check out

1. Securing docker containers
2. Profiling heap and CPU (pprof) in detail

## Interesting things I learnt while working on this

1. Docker build optimixation - [blog](https://blog.pranshu-raj.me/posts/optimizing-docker-builds/)
2. Grafana dashboard config for persistence across builds (blog on this soon)
3. Mutexes, race conditions (got to see these happen in my code)
4. Go concurrency - goroutines and channels (and actors)

## Things I need to learn more about
1. Testing
2. Go context
3. Profiling
4. Go garbage collector
