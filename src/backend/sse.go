package backend

import (
	"encoding/json"
	"fmt"
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/redisclient"
	"time"

	"github.com/gin-gonic/gin"
)

func StreamLeaderboard(c *gin.Context) {
	metrics.ConcurrentClients.Inc()
	defer metrics.ConcurrentClients.Dec()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	metrics.ActiveSSEConnections.Inc()
	defer metrics.ActiveSSEConnections.Dec()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var lastData []byte

	for {
		select {
		case <-ticker.C:
			results, err := redisclient.GetTopNPlayers(c, "leaderboard", int64(config.AppConfig.Leaderboard.TopPlayersLimit))
			if err != nil {
				metrics.RedisOperationErrors.WithLabelValues("get_top_players").Inc()
			}

			jsonStart :=time.Now()

			data, err := json.Marshal(results)
			metrics.JSONMarshalDuration.Observe(time.Since(jsonStart).Seconds())
    
			if err != nil {
				metrics.JSONErrors.WithLabelValues("marshal").Inc()
				return
			}


			if !jsonEqual(data, lastData) {
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				c.Writer.Flush()
				metrics.SSEMessagesSent.Inc()
				lastData = data
			}

		case <-c.Request.Context().Done():
			metrics.DroppedSSEConnections.Inc()
			return
		}
	}
}

func jsonEqual(a, b []byte) bool {
	return string(a) == string(b)
}
