package models

import (
	"fmt"
)

type GameResult struct {
	GameID    string `json:"game_id" binding:"required"`
	ServerID  string `json:"server_id" binding:"required"`
	Player1ID string `json:"player1_id" binding:"required"`
	Player2ID string `json:"player2_id" binding:"required"`
	Result    int    `json:"result" binding:"gte=0,lte=2"`
}

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

func (g *GameResult) Validate() error {
	if g.GameID == "" {
		return &ValidationError{Field:"game_id", Reason:"empty string not allowed"}
	}
	if g.Player1ID == g.Player2ID {
		return &ValidationError{Field: "player_id", Reason: "player ids cannot be the same"}
	}
	if g.Result < 0 || g.Result > 2 {
		return &ValidationError{Field: "result", Reason: "result must be an integer between 0 and 2 inclusive"}
	}
	if g.ServerID == "" {
		return &ValidationError{Field: "server_id", Reason: "empty string not allowed"}
	}
	return nil
}
