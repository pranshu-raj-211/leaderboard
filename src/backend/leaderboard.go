package backend

import (
	"encoding/json"
	"leaderboard/src/config"
	"leaderboard/src/redisclient"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {
	results, err := redisclient.GetTopNPlayers(c, "leaderboard", 0)
	if err != nil {
		config.Error("Could not fetch leaderboard from Redis", map[string]any{"Error": err})
	}

	data, _ := json.Marshal(results)
	c.Writer.Write(data)
}
