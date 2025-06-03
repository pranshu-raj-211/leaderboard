package models

import (
	"leaderboard/src/config"
)

type GameResult struct {
	GameID    string `json:"game_id" binding:"required"`
	ServerID  string `json:"server_id" binding:"required"`
	Player1ID string `json:"player1_id" binding:"required"`
	Player2ID string `json:"player2_id" binding:"required"`
	Result    int    `json:"result" binding:"gte=0,lte=2"`
}

func (g *GameResult) Validate() error {
	if g.GameID == "" {
		return config.Error("game_id cannot be empty", map[string]any{"GameID": g.GameID})
	}
	if g.Player1ID == g.Player2ID {
		return config.Error("Player IDs cannot be the same", map[string]any{"Player1ID": g.Player1ID, "Player2ID": g.Player2ID, "GameID": g.GameID})
	}
	if g.Result < 0 || g.Result > 2 {
		return config.Error("Game result must be between 0 and 2", map[string]any{"GameID": g.GameID, "Result": g.Result})
	}
	return nil
}
