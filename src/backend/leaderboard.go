package backend

import (
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(store interfaces.LeaderboardStore, limit int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Use limit=0 if we want to get the whole leaderboard (0, -1)
		results, err := store.GetTopNPlayers(ctx, limit)
		if err != nil {
			config.Error("Could not fetch leaderboard from Redis", map[string]any{"Error": err})
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch leaderboard"})
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}
