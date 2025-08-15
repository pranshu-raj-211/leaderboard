package redisclient

import (
	"context"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	maxRetries := config.AppConfig.Redis.MaxRetries
	for i := 0; i < maxRetries; i++ {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.AppConfig.Redis.Address,
			Password: config.AppConfig.Redis.Password,
			DB:       config.AppConfig.Redis.DB,
		})

		_, err := redisClient.Ping(context.Background()).Result()
		// keep retrying until redis is connected, fatal otherwise
		if err == nil {
			return
		}
		config.Error("Failed to connect to redis client, retrying", map[string]any{"Error": err})
		time.Sleep(2 * time.Second)
	}
	config.Fatal("Could not connect to redis client after max retries", map[string]any{"Tried":maxRetries})
}

// The function also observes RedisLatency and LeaderboardUpdateDuration metrics after performing the update.
func UpdateLeaderboard(ctx context.Context, player1ID, player2ID string, result int) error {
	start := time.Now()
	updateStart := time.Now()
	switch result {
	case 0:
		//player 1 wins
		redisClient.ZIncrBy(ctx, "leaderboard", 1.0, player1ID)
	case 1:
		//player 2 wins
		redisClient.ZIncrBy(ctx, "leaderboard", 1.0, player2ID)
	case 2:
		//draw due to any reason
		redisClient.ZIncrBy(ctx, "leaderboard", 0.5, player1ID)
		redisClient.ZIncrBy(ctx, "leaderboard", 0.5, player2ID)
	default:
		return config.Error("Invalid game result, did not update leaderboard",
			map[string]any{
				"player1ID": player1ID,
				"player2ID": player2ID,
				"result":    result,
			})
	}
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	metrics.LeaderboardUpdateDuration.Observe(time.Since(updateStart).Seconds())
	return nil
}

// GetTopNPlayers returns up to n entries from the sorted set stored at key ordered by highest score first.
// It calls Redis ZRevRangeWithScores to fetch the top N members with their scores and returns the resulting []redis.Z.
// On success, returns top N scores, sorted in reverse order (highest first). On failure, records error, returns.
func GetTopNPlayers(ctx context.Context, key string, n int64) ([]redis.Z, error) {
	start := time.Now()

	scores, err := redisClient.ZRevRangeWithScores(ctx, key, 0, n-1).Result()

	if err != nil {
		metrics.RedisLatency.Observe(time.Since(start).Seconds())
		return nil, config.Error("Failed to fetch top n players", map[string]any{})
	}
	metrics.RedisPayloadSize.Observe(float64(len(scores)))
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	return scores, nil
}

// GetPlayerScore returns the leaderboard rank and score for the given player ID in the specified sorted set key.
//
// The function queries Redis for the member's rank and score and observes the Redis latency metric before returning.
// If the player is not present, Redis returns redis.Nil; in that case the error will be redis.Nil and the returned
// rank and score will be zero. Any other Redis error is returned to the caller as well.
func GetPlayerScore(ctx context.Context, key string, playerID string) (int64, float64, error) {
	start := time.Now()

	player_info, err := redisClient.ZRankWithScore(ctx, key, playerID).Result()
	if err != nil {
		metrics.RedisLatency.Observe(time.Since(start).Seconds())
		return 0, 0, config.Error("Something went wrong while getting player stats", map[string]any{"player_id": playerID, "Error": err})
	}
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	rank := player_info.Rank
	score := player_info.Score
	return rank, score, nil
}
