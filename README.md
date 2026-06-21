# Real-time Game Leaderboard

Blogs on this project - [link](https://blog.pranshu-raj.me/tags/leaderboard/)

Some recent posts on this:
- [Breaking the 28k connection barrier](https://blog.pranshu-raj.me/posts/scaling-sse-1m-connections)
- [28k+ conns, zero messages](https://blog.pranshu-raj.me/posts/implementing-correct-fanout)
- [Optimizing docker image builds](https://blog.pranshu-raj.me/posts/optimizing-docker-builds)
- [Backpressure](https://blog.pranshu-raj.me/posts/understanding-backpressure/)


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

## Prerequisites to run

- Docker (and Docker Compose)

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

Currently using Go `1.23.4`

(Optional) Configure the application in `config.yaml`:

To checkout the Grafana dashboard, run the app and go to `http://localhost:3000`, login with the default Grafana username and password (admin).

## API Endpoints

### Submit Game Results
see `src/backend/submit_game.go`

Submits a game result JSON to the app, to simulate game servers sending game results.

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

### Register Player 
see `src/backend/player.go`

Maps a player ID to a human-readable display name. The name is stored in a Redis hash and is returned alongside the player ID on all reads (this wasn't really needed from a backend pov, but good to see from a user's pov).

```http
POST /players
Content-Type: application/json

{
  "player_id": "player1",
  "name": "SwiftFalcon"
}
```

### Get Leaderboard
see `src/backend/leaderboard.go`

Gets a snapshot of the leaderboard at that instant.

```http
GET /leaderboard
```

```json
[
  {"player_id": "p_0009", "name": "BlazingFalcon9", "score": 5},
  {"player_id": "p_0013", "name": "BraveTiger13", "score": 4.5}
]
```

### Get Player Stats
see `src/backend/player.go`

See the rank, score, and name of a player by their id.

```http
GET /player/:id/stats
```

```json
{"player_id": "p_0000", "name": "IronFalcon0", "rank": 8, "score": 2.5}
```

### Stream Real-time Updates
see `src/backend/sse.go`

Stream leaderboard updates periodically to the client connection.

```http
GET /stream-leaderboard
```

## Match Simulator (write traffic)
see `src/simulator/simulator.go`

Previously the system only had a reader path (the SSE broadcaster polling Redis) and nothing wrote game results, so the leaderboard never updated. This is the most important feature for testing, on startup it seeds a pool of random players, then on every tick plays a random number of concurrent matches between random pairs of players, writing the outcome to the write endpoint.

Configured under `simulator:` in `config.yaml`:

```yaml
simulator:
  enabled: true
  num_players: 50
  tick_interval_millis: 2000
  max_concurrent_matches: 5
  seed: 42                      # to satisfy the PRNG gods (deterministic testing)
```

### Prometheus Metrics
see `src/metrics/metrics.go`

For prometheus connection.

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
|   ├── interfaces/  # interfaces - store and broadcaster
│   └── main.go      # Application entry point
|── tests/
|   |── unit/        # unit testing scripts
├── config.yaml      # Configuration file
├── go.mod           # Go module file
└── README.md        # This file
```

## Benchmarking
To load test the application, first spin up the app using docker compose.

`docker compose up --build`

Then in a separate shell, after the compose app is up and running, run the client docker containers using the command

```bash
docker run -d --name client1 --network=leaderboard_internal sseclient:0.2 -ip=leaderboard -conn=20000
```

The image needs to be built before though, from the Dockerfile in the `client` directory.

---

The server and client containers on a single machine reaches up to 150k concurrent SSE connections, which is documented in [this blog](https://blog.pranshu-raj.me/posts/scaling-sse-1m-connections/). At this point we hit memory limits (testing on a Linux laptop with 8GB RAM, so results may not be consistent each time as with laptops). This was tested without game servers submitting results, which isn't a very realistic benchmark.

[Post with Grafana dashboard - 28231 active SSE conns](https://x.com/seigino99707047/status/1955225344744583581)

[Post with game submissions dashboard](https://x.com/seigino99707047/status/1955246776035721718)

Combined streamer and submitter scripts have not been tested beyond this limit, however the streamer was tested in standalone mode, which went to 28,232 connections (default max outbound ports in Linux).

The streaming script is based off of [Eran Yanay's repo for his Gophercon talk](https://github.com/eranyanay/1m-go-websockets).

## Ideas to check out

1. Securing docker containers
2. Profiling heap and CPU (pprof) in detail

## Interesting things I learnt while working on this

1. Docker build optimixation - [blog](https://blog.pranshu-raj.me/posts/optimizing-docker-builds/)
2. Grafana dashboard config for persistence across builds (done - blog on this soon)
3. Mutexes, race conditions (got to see these happen in my code)
4. Go concurrency - goroutines and channels (and actors)

## Things I need to learn more about
1. Testing
2. Go context
3. Profiling
4. Go garbage collector
