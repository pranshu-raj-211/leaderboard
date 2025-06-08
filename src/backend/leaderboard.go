package backend

import (
	"encoding/json"
	"leaderboard/src/config"
	"leaderboard/src/redisclient"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(c *gin.Context) {
	// Using n=0 since we want to get the whole leaderbord (0, -1)
	results, err := redisclient.GetTopNPlayers(c, "leaderboard", 0)
	if err != nil {
		config.Error("Could not fetch leaderboard from Redis", map[string]any{"Error": err})
		c.JSON(500, gin.H{"error": "could not fetch leaderboard"})
		return
	}

	data, _ := json.Marshal(results)
	c.Writer.Write(data)
}
