package main

import (
	"fmt"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	config.InitLogger()
	if err := config.LoadConfig("config.yaml"); err != nil {
		config.Fatal("Failed to load config")
	}
	config.Info("Starting server", zap.String("redis_addr", config.AppConfig.Redis.Address), zap.Int("port", config.AppConfig.Server.Port))

	redisclient.InitRedis()
	metrics.InitMetrics()

	r := gin.Default()
	r.POST("/submit-game", backend.SubmitGameResults)
	r.GET("/stream-leaderboard", backend.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	address := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	r.Run(address)
}
