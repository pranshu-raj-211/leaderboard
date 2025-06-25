package backend

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type LeaderboardUpdate struct {
	Data []byte
	Hash [32]byte
}

type LeaderboardBroadcaster struct {
	// channel for all SSE conns to listen to get lb updates
	broadcastChan chan LeaderboardUpdate

	ctx    context.Context
	cancel context.CancelFunc

	// enable goroutines to run before moving on, to be used for locks
	wg sync.WaitGroup
}

func CreateLeaderboardBroadcaster() *LeaderboardBroadcaster {
	ctx, cancel := context.WithCancel(context.Background())

	lb := &LeaderboardBroadcaster{
		// make the channel buffered - clients may be slow, messages can pile up
		broadcastChan: make(chan LeaderboardUpdate, config.AppConfig.Server.BroadcastBufferSize),
		ctx:           ctx,
		cancel:        cancel,
	}

	lb.wg.Add(1)
	go lb.detectLeaderboardChanges()
	return lb
}

func (lb *LeaderboardBroadcaster) StopBroadcast() {
	lb.cancel()
	lb.wg.Wait()
	close(lb.broadcastChan)
}

func (lb *LeaderboardBroadcaster) GetBroadcastChannel() <-chan LeaderboardUpdate {
	return lb.broadcastChan
}

// package level var
var broadcaster *LeaderboardBroadcaster

func SetBroadcaster(b *LeaderboardBroadcaster) {
	broadcaster = b
}

func (lb *LeaderboardBroadcaster) detectLeaderboardChanges() {
	defer lb.wg.Done()

	// TODO: check this time conversion
	ticker := time.NewTicker(time.Duration(config.AppConfig.Server.PollingIntervalSeconds) * time.Second)
	defer ticker.Stop()

	var lastHash [32]byte

	for {
		select {
		case <-ticker.C:
			results, err := redisclient.GetTopNPlayers(lb.ctx, "leaderboard", int64(config.AppConfig.Leaderboard.TopPlayersLimit))
			if err != nil {
				metrics.RedisOperationErrors.WithLabelValues("get_top_players").Inc()
				config.Error("Failed to fetch leaderboard from Redis.", map[string]any{"Error": err, "source": "/stream-leaderboard"})
				continue
			}

			resultString := fmt.Sprintf("%+v", results)
			currentHash := sha256.Sum256([]byte(resultString))

			if currentHash != lastHash {
				lastHash = currentHash

				jsonStart := time.Now()
				jsonData, err := json.Marshal(results)
				if err != nil {
					config.Error("JSON marshaling error",
						map[string]any{"Error": err, "source": "/stream-leaderboard", "results": results})
					metrics.JSONErrors.WithLabelValues("marshal").Inc()
					continue
				}
				metrics.JSONMarshalDuration.Observe(time.Since(jsonStart).Seconds())

				update := LeaderboardUpdate{
					Data: jsonData,
					Hash: currentHash,
				}

				// non blocking send
				select {
				case lb.broadcastChan <- update:
				default:
				}
			}
		case <-lb.ctx.Done():
			return
		}
	}
}

func StreamLeaderboard(c *gin.Context) {
	metrics.ConcurrentClients.Inc()
	defer metrics.ConcurrentClients.Dec()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	metrics.ActiveSSEConnections.Inc()
	config.Info("New SSE conn", map[string]any{"Num active clients": metrics.ActiveSSEConnections})
	defer metrics.ActiveSSEConnections.Dec()

	broadcastChan := broadcaster.GetBroadcastChannel()

	for {
		select {
		case update, ok := <-broadcastChan:
			if !ok {
				// channel closed
				return
			}
			fmt.Fprintf(c.Writer, "data: %s\n\n", update.Data)
			c.Writer.Flush()
			metrics.SSEMessagesSent.Inc()

		case <-c.Request.Context().Done():
			metrics.DroppedSSEConnections.Inc()
			config.Info("Closed SSE conn", map[string]any{"open": metrics.ActiveSSEConnections})
			return
		}
	}
}
