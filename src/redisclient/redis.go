package redisclient

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to redis client: %v", err))
	}
}

func UpdateLeaderboard(ctx context.Context, player1ID, player2ID string, result int) error {
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
		return fmt.Errorf("invalid result code")
	}
	return nil
}

func GetTopNPlayers(ctx context.Context, key string, n int64) ([]redis.Z, error) {
	scores, err := redisClient.ZRevRangeWithScores(ctx, key, 0, n-1).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch leaderboard")
	}
	return scores, nil
}
