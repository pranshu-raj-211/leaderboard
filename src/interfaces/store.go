package interfaces

import "context"

type LeaderboardEntry struct{
	PlayerID string
	Score float64
}

type LeaderboardStore interface {
	UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error
	GetTopNPlayers(ctx context.Context, limit int64) ([]LeaderboardEntry, error)
	GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error)
}