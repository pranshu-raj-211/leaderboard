package main

import (
	"context"
	"fmt"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	cfg := &backend.BroadcasterConfig{
		BroadcastBufferSize:      config.AppConfig.Server.BroadcastBufferSize,
		PollingIntervalSeconds:   config.AppConfig.Leaderboard.UpdateIntervalSecs,
		TopPlayersLimit:          config.AppConfig.Leaderboard.TopPlayersLimit,
		HeartbeatIntervalSeconds: config.AppConfig.Server.HeartbeatIntervalSeconds,
	}

	metrics.InitMetrics()
	// naming will need changes, future postgres dep will add a postgres client
	client, err := redisclient.NewRedisClient(*config.AppConfig, config.GetLogger())
	if err != nil {
		config.Fatal("Could not connect to redis, shutting down", map[string]any{"err": err})
	}

	store := redisclient.CreateRedisLeaderboard(client)

	lb, err := backend.CreateLeaderboardBroadcaster(store, cfg)
	if err != nil {
		config.Fatal("incorrect config passed to leaderboard constructor", map[string]any{"err": err})
	}
	defer lb.StopBroadcast()

	r := gin.Default()

	r.Use(metrics.MetricsMiddleware())

	r.POST("/submit-game", backend.SubmitGameResults(store))
	r.GET("/stream-leaderboard", lb.StreamLeaderboard)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/player/:id/stats", backend.GetPlayerResults(store))
	r.GET("/leaderboard", backend.GetLeaderboard(store, int64(config.AppConfig.Leaderboard.TopPlayersLimit)))

	address := fmt.Sprintf("%s:%d", config.AppConfig.Server.Host, config.AppConfig.Server.Port)

	srv := &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Fatal("listen error", map[string]any{"err": err})
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	config.Info("Shutdown signal received, gracefully shutting down server...", map[string]any{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.AppConfig.Server.GracefulShutdownTimeoutSeconds))
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		config.Error("Server forced to shutdown", map[string]any{"err": err})
	}

	lb.StopBroadcast()
	config.Info("Server exiting", map[string]any{})
}
