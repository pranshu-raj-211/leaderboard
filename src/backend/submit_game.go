package backend

import (
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"leaderboard/src/metrics"
	"leaderboard/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubmitGameResults accepts a JSON-encoded game result from a game server, updates the leaderboard in Redis,
// increments the submission metric, and responds with an HTTP status.
//
// On invalid JSON the handler responds with HTTP 400 and an error message. If updating the leaderboard fails
// it responds with HTTP 400 and the underlying error message. On success it responds with HTTP 200.
func SubmitGameResults(store interfaces.LeaderboardStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var game models.GameResult
		if err := c.ShouldBindJSON(&game); err != nil {
			config.Error("Invalid JSON received from game server", map[string]any{"Error": err})
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json for game result"})
			return
		}

		if err := game.Validate(); err != nil {
			config.Error("Game object validation error", map[string]any{"Error": err})
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game result"})
			return
		}

		if err := store.UpdateLeaderboard(c.Request.Context(), game.Player1ID, game.Player2ID, game.Result); err != nil {
			config.Error("Could not update leaderboard", map[string]any{"Error": err, "GameID": game.GameID})
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		metrics.GameSubmissions.Inc()
		c.JSON(http.StatusOK, gin.H{"status": "Leaderboard updated"})
	}
}
