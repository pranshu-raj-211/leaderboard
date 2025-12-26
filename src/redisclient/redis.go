package redisclient

import (
	"context"
	"errors"
	"fmt"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"leaderboard/src/metrics"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisLeaderboard struct {
	client *redis.Client
}

func NewRedisClient(cfg config.Config, logger *zap.Logger) (*redis.Client, error) {
	for i := 0; i < cfg.Redis.MaxRetries; i++ {
		client := redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := client.Ping(ctx).Result()
		cancel()

		if err == nil {
			return client, nil
		}
		logger.Warn(
			"redis connection failed, retrying",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to redis after %d retries, closing", cfg.Redis.MaxRetries)
}

func CreateRedisLeaderboard(client *redis.Client) *RedisLeaderboard {
	return &RedisLeaderboard{client: client}
}

// The function also observes RedisLatency and LeaderboardUpdateDuration metrics after performing the update.
func (store *RedisLeaderboard) UpdateLeaderboard(ctx context.Context, player1ID, player2ID string, result int) error {
	start := time.Now()

	var err error

	switch result {
	case 0:
		//player 1 wins
		err = store.client.ZIncrBy(ctx, "leaderboard", 1.0, player1ID).Err()
	case 1:
		//player 2 wins
		err = store.client.ZIncrBy(ctx, "leaderboard", 1.0, player2ID).Err()
	case 2:
		//draw due to any reason
		pipe := store.client.Pipeline()
		pipe.ZIncrBy(ctx, "leaderboard", 0.5, player1ID)
		pipe.ZIncrBy(ctx, "leaderboard", 0.5, player2ID)
		_, err = pipe.Exec(ctx)
	default:
		// TODO: return typed error instead of cfg.error, repeat across other places as well
		return config.Error("Invalid game result, did not update leaderboard",
			map[string]any{
				"player1ID": player1ID,
				"player2ID": player2ID,
				"result":    result,
			})
	}
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	metrics.LeaderboardUpdateDuration.Observe(time.Since(start).Seconds())

	if err != nil {
		return config.Error("Failed to update leaderboard", map[string]any{"Error": err, "player1ID": player1ID, "player2ID": player2ID, "Result": result})
	}

	return nil
}

// GetTopNPlayers returns up to n entries from the sorted set stored at key ordered by highest score first.
// It calls Redis ZRevRangeWithScores to fetch the top N members with their scores and returns the resulting []redis.Z.
// On success, returns top N scores, sorted in reverse order (highest first). On failure, records error, returns.
func (store *RedisLeaderboard) GetTopNPlayers(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
	start := time.Now()

	zs, err := store.client.ZRevRangeWithScores(ctx, "leaderboard", 0, limit-1).Result()
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	if err != nil {
		return nil, config.Error("Failed to fetch top n players", map[string]any{"Error": err})
	}
	entries := make([]interfaces.LeaderboardEntry, 0, len(zs))
	for _, z := range zs {
		playerID, ok := z.Member.(string)
		if !ok {
			return nil, fmt.Errorf("invalid redis member type")
		}

		entries = append(entries, interfaces.LeaderboardEntry{
			PlayerID: playerID,
			Score:    z.Score,
		})
	}

	return entries, nil
}

// GetPlayerScore returns the leaderboard rank and score for the given player ID in the specified sorted set key.
//
// The function queries Redis for the member's rank and score and observes the Redis latency metric before returning.
// If the player is not present, Redis returns redis.Nil; in that case the 0,0, nil is returned.
// Any other Redis error is returned to the caller.
func (store *RedisLeaderboard) GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error) {
	start := time.Now()

	playerInfo, err := store.client.ZRankWithScore(ctx, "leaderboard", playerID).Result()
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	if errors.Is(err, redis.Nil) {
		//player not found
		config.Info("Player not found", map[string]any{"Player ID": playerID, "source": "GetPlayerScore"})
		return 0, 0, nil
	}
	if err != nil {
		return 0, 0, config.Error("Something went wrong while getting player stats", map[string]any{"player_id": playerID, "Error": err})
	}

	return playerInfo.Rank, playerInfo.Score, nil
}
