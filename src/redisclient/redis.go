package redisclient

import (
	"context"
	"leaderboard/src/config"

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
		config.Error("Failed to connect to redis client", map[string]any{"Error": err})
		panic("Trying recovery for redis client connection")
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
		return config.Error("Invalid game result, did not update leaderboard",
			map[string]any{
				"player1ID": player1ID,
				"player2ID": player2ID,
				"result":    result,
			})
	}
	return nil
}

func GetTopNPlayers(ctx context.Context, key string, n int64) ([]redis.Z, error) {
	scores, err := redisClient.ZRevRangeWithScores(ctx, key, 0, n-1).Result()

	if err != nil {
		return nil, config.Error("Failed to fetch top n players", map[string]any{})
	}
	return scores, nil
}

func GetPlayerScore(ctx context.Context, key string, playerID string) (int64, float64, error) {
	player_info, err := redisClient.ZRankWithScore(ctx, key, playerID).Result()
	if err != nil {
		config.Error("Something went wrong while getting player stats", map[string]any{"player_id": playerID, "Error": err})
	}
	rank := player_info.Rank
	score := player_info.Score
	// TODO: checkout result of this operation, split to get rank and score
	return rank, score, err
}
