package main

import (
	"fmt"
	"leaderboard/src/backend"
	"leaderboard/src/redisclient"
	"leaderboard/src/metrics"
	"leaderboard/src/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)



func main() {
	if err := config.LoadConfig("config.yaml"); err != nil {
        panic(err)
    }
    fmt.Println(config.AppConfig.Redis.Address)

	redisclient.InitRedis()
	metrics.InitMetrics()

	r := gin.Default()
	r.POST("/submit-game", backend.SubmitGameResults)
	r.GET("/stream-leaderboard", backend.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run(":8080")
}
