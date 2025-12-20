package main

import (
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

	var broadcaster broadcaster = backend.CreateLeaderboardBroadcaster()
	defer broadcaster.StopBroadcast()

	r := gin.Default()

	r.Use(metrics.MetricsMiddleware())

	r.POST("/submit-game", backend.SubmitGameResults)
	r.GET("/stream-leaderboard", broadcaster.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/player/:id/stats", backend.GetPlayerResults)
	r.GET("/leaderboard", backend.GetLeaderboard)

	address := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	r.Run(address)
}
