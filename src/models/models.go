package models

import "fmt"

type GameResult struct {
	GameID    string `json:"game_id" binding:"required"`
	ServerID  string `json:"server_id" binding:"required"`
	Player1ID string `json:"player1_id" binding:"required"`
	Player2ID string `json:"player2_id" binding:"required"`
	Result    int    `json:"result" binding:"gte=0,lte=2"`
}

func (g *GameResult) Validate() error {
	if g.GameID == "" {
		return fmt.Errorf("game_id cannot be empty")
	}
	if g.Player1ID == g.Player2ID {
		return fmt.Errorf("player1_id and player2_id cannot be the same")
	}
	if g.Result < 0 || g.Result > 2 {
		return fmt.Errorf("result must be between 0 and 2")
	}
	return nil
}
