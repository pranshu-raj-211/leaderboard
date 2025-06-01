package models

type GameResult struct {
	GameID    string `json:"game_id"`
	ServerID  string `json:"server_id"`
	Player1ID string `json:"player1_id"`
	Player2ID string `json:"player2_id"`
	Result    int    `json:"result"`
}
