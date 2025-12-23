package main

import (
	"context"
	"fmt"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type broadcaster interface {
	StopBroadcast()
	StreamLeaderboard(*gin.Context)
}

type RedisLeaderboard struct {}

func (RedisLeaderboard) UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error {
	return redisclient.UpdateLeaderboard(ctx, player1ID, player2ID, result)
}

func main() {
	// TODO: use config loading before logger init (use config vars in InitLogger)
	config.InitLogger()
	if err := config.LoadConfig("config.yaml"); err != nil {
		config.Fatal("Failed to load config", map[string]any{"err": err})
	}
	config.Info("Starting server", map[string]any{
		"Redis address": config.AppConfig.Redis.Address,
		"Host":          config.AppConfig.Server.Host,
		"Port":          config.AppConfig.Server.Port,
	})

	redisclient.InitRedis()
	metrics.InitMetrics()

	// TODO: checkout why to prefer broadcast interface here
	var lb broadcaster = backend.CreateLeaderboardBroadcaster()
	defer lb.StopBroadcast()

	store := RedisLeaderboard{}

	r := gin.Default()

	r.Use(metrics.MetricsMiddleware())

	r.POST("/submit-game", backend.SubmitGameResults(store))
	r.GET("/stream-leaderboard", lb.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/player/:id/stats", backend.GetPlayerResults)
	r.GET("/leaderboard", backend.GetLeaderboard)

	address := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	r.Run(address)
}
