package backend

import (
	"encoding/json"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"

	"github.com/gin-gonic/gin"
)

func GetLeaderboard(store interfaces.LeaderboardStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Use n=0 if we want to get the whole leaderbord (0, -1)
		results, err := store.GetTopNPlayers(ctx, 10)
		if err != nil {
			config.Error("Could not fetch leaderboard from Redis", map[string]any{"Error": err})
			ctx.JSON(500, gin.H{"error": "could not fetch leaderboard"})
			return
		}

		data, _ := json.Marshal(results)
		ctx.Writer.Write(data)
	}
}
