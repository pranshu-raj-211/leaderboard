package backend

import (
	"leaderboard/src/config"
	"leaderboard/src/metrics"
	"leaderboard/src/models"
	"leaderboard/src/redisclient"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitGameResults(c *gin.Context) {
	metrics.GameSubmissions.Inc()
	var game models.GameResult
	if err := c.ShouldBindJSON(&game); err != nil {
		config.Error("Invalid JSON received from game server", map[string]any{"Error": err})
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json for game result"})
		return
	}

	if err := redisclient.UpdateLeaderboard(c.Request.Context(), game.Player1ID, game.Player2ID, game.Result); err != nil {
		config.Error("Could not update leaderboard", map[string]any{"Error": err, "GameID": game.GameID})
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Leaderboard updated"})
}
