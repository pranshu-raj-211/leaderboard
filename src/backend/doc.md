The backend package implements HTTP handlers for the leaderboard service's REST API endpoints.

## Handlers

- `SubmitGameResults`: Accepts game results and updates player scores
- `StreamLeaderboard`: Implements SSE for real-time leaderboard updates
- `GetLeaderboard`: Returns current leaderboard standings
- `GetPlayerResults`: Fetches individual player statistics

## Architecture

- Uses Gin framework for routing
- Integrates with Redis for data storage
- Implements SSE for real-time updates
- Includes Prometheus metrics