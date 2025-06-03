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

func main() {
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

	r := gin.Default()

	r.Use(metrics.MetricsMiddleware())

	r.POST("/submit-game", backend.SubmitGameResults)
	r.GET("/stream-leaderboard", backend.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	address := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	r.Run(address)
}
