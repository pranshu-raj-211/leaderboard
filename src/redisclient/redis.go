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
		if err == nil {
			return
		}
		config.Error("Failed to connect to redis client, retrying", map[string]any{"Error": err})
		time.Sleep(2 * time.Second)
	}
	config.Fatal("Could not connect to redis client after 10 retries", map[string]any{})
}

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

func GetTopNPlayers(ctx context.Context, key string, n int64) ([]redis.Z, error) {
	start := time.Now()

	scores, err := redisClient.ZRevRangeWithScores(ctx, key, 0, n-1).Result()

	metrics.RedisPayloadSize.Observe(float64(len(scores)))
    metrics.RedisLatency.Observe(time.Since(start).Seconds())
    

	if err == redis.Nil {
		return nil, config.Error("Failed to fetch top n players", map[string]any{})
	}
	return scores, nil
}

func GetPlayerScore(ctx context.Context, key string, playerID string) (int64, float64, error) {
	start := time.Now()

	player_info, err := redisClient.ZRankWithScore(ctx, key, playerID).Result()
	if err == redis.Nil {
		config.Error("Something went wrong while getting player stats", map[string]any{"player_id": playerID, "Error": err})
	}
	metrics.RedisLatency.Observe(time.Since(start).Seconds())
	rank := player_info.Rank
	score := player_info.Score
	return rank, score, err
}
