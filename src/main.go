package main

import (
	"fmt"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

	metrics.InitMetrics()
	// naming will need changes, future postgres dep will add a postgres client
	client, err := redisclient.NewRedisClient(*config.AppConfig, config.GetLogger())
	if err!= nil{
		config.Fatal("Could not connect to redis, shutting down", map[string]any{"err":err})
	}

	store := redisclient.CreateRedisLeaderboard(client)

	// TODO: checkout why to prefer broadcast interface here
	var lb interfaces.Broadcaster = backend.CreateLeaderboardBroadcaster(store)
	defer lb.StopBroadcast()

	r := gin.Default()

	r.Use(metrics.MetricsMiddleware())

	r.POST("/submit-game", backend.SubmitGameResults(store))
	r.GET("/stream-leaderboard", lb.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/player/:id/stats", backend.GetPlayerResults(store))
	r.GET("/leaderboard", backend.GetLeaderboard(store))

	address := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	r.Run(address)
}
