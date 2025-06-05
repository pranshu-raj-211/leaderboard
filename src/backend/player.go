package backend

import (
	"leaderboard/src/config"
	"leaderboard/src/redisclient"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayerResults(c *gin.Context) {
	playerID := c.Param("id")
	if playerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "player id required to fetch stats."})
		return
	}
	rank, score, err := redisclient.GetPlayerScore(c, "leaderboard", playerID)
	if err != nil {
		config.Error("Error getting player score", map[string]any{"Error": err})
	}
	config.Info("Results from player stats api", map[string]any{"id": playerID, "rank": rank, "score": score})
	c.JSON(http.StatusOK, gin.H{"rank": rank, "score": score})
}
