package backend

import (
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlayerResults(store interfaces.LeaderboardStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		playerID := ctx.Param("id")
		if playerID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "player id required to fetch stats."})
			return
		}
		rank, score, err := store.GetPlayerScore(ctx, playerID)
		if err != nil {
			config.Error("Error getting player score", map[string]any{"Error": err})
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch player stats"})
			return
		}
		config.Info("Results from player stats api", map[string]any{"id": playerID, "rank": rank, "score": score})
		ctx.JSON(http.StatusOK, gin.H{"rank": rank, "score": score})
	}
}
