package backend

import (
	"encoding/json"
	"fmt"
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

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			results, err := redisclient.GetTopNPlayers(c, "leaderboard", 10)
			if err != nil {
				continue
			}

			data, _ := json.Marshal(results)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()

		case <-c.Request.Context().Done():
			return
		}
	}
}
