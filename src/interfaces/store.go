package interfaces

import "context"

type LeaderboardEntry struct {
	PlayerID string  `json:"player_id"`
	Name     string  `json:"name"`
	Score    float64 `json:"score"`
}

type LeaderboardStore interface {
	UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error
	GetTopNPlayers(ctx context.Context, limit int64) ([]LeaderboardEntry, error)
	GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error)

	SetPlayerName(ctx context.Context, playerID string, name string) error
	GetPlayerName(ctx context.Context, playerID string) (string, error)
	// GetPlayerNames batch-resolves names for the given IDs. IDs without a stored name are simply absent from the returned map.
	GetPlayerNames(ctx context.Context, playerIDs []string) (map[string]string, error)
}
