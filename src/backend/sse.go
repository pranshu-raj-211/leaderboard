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
				continue
			}

			data, _ := json.Marshal(results)
			if !jsonEqual(data, lastData) {
				fmt.Fprintf(c.Writer, "data: %s\n\n", data)
				c.Writer.Flush()
				lastData = data
			}

		case <-c.Request.Context().Done():
			return
		}
	}
}

func jsonEqual(a, b []byte) bool {
	return string(a) == string(b)
}
