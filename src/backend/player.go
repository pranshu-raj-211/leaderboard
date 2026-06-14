package backend

import (
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"leaderboard/src/models"
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
		name, err := store.GetPlayerName(ctx, playerID)
		if err != nil {
			config.Error("Error getting player name", map[string]any{"Error": err})
		}
		config.Info("Results from player stats api", map[string]any{"id": playerID, "name": name, "rank": rank, "score": score})
		ctx.JSON(http.StatusOK, gin.H{"player_id": playerID, "name": name, "rank": rank, "score": score})
	}
}

//Accepts a JSON {player_id, name} and stores the mapping so that reads (leaderboard, stats, SSE) can return human-readable names alongside IDs.
func RegisterPlayer(store interfaces.LeaderboardStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var player models.Player
		if err := ctx.ShouldBindJSON(&player); err != nil {
			config.Error("Invalid JSON received for player registration", map[string]any{"Error": err})
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid json for player"})
			return
		}
		if err := player.Validate(); err != nil {
			config.Error("Player validation error", map[string]any{"Error": err})
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := store.SetPlayerName(ctx, player.PlayerID, player.Name); err != nil {
			config.Error("Could not register player", map[string]any{"Error": err, "player_id": player.PlayerID})
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not register player"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "player registered", "player_id": player.PlayerID, "name": player.Name})
	}
}
