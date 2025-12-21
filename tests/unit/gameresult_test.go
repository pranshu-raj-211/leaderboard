package unit_test

import (
	"leaderboard/src/models"
	"testing"
)

type testCase struct {
	name         string
	input        models.GameResult
	wantErr      bool
	wantErrField string
}

var testCases = []testCase{
	{
		name:         "Valid Player1 win",
		input:        models.GameResult{GameID: "123", Player1ID: "1", Player2ID: "2", ServerID: "2", Result: 2},
		wantErr:      false,
		wantErrField: "",
	},
	{
		name: "valid result",
		input: models.GameResult{
			GameID:    "g1",
			ServerID:  "s2",
			Player1ID: "p1",
			Player2ID: "p2",
			Result:    2,
		},
		wantErr:      false,
		wantErrField: "",
	},
	{
		name: "same player ids",
		input: models.GameResult{
			GameID:    "g1",
			ServerID:  "s2",
			Player1ID: "p1",
			Player2ID: "p1",
			Result:    1,
		},
		wantErr:      true,
		wantErrField: "player_id",
	},
	{
		name: "result too high",
		input: models.GameResult{
			GameID:    "g1",
			Player1ID: "p1",
			Player2ID: "p2",
			Result:    3,
		},
		wantErr:      true,
		wantErrField: "result",
	},
	{
		name: "Empty Server ID",
		input: models.GameResult{
			GameID:    "g1",
			ServerID:  "",
			Player1ID: "p1",
			Player2ID: "p2",
			Result:    2,
		},
		wantErr:      true,
		wantErrField: "server_id",
	},
}

func TestGameResult(t *testing.T) {
	for _, sample := range testCases {
		t.Run(sample.name, func(t *testing.T) {
			got := sample.input.Validate()

			if got == nil && sample.wantErr {
				t.Errorf("expected an error, got none")
			}
			if got != nil && !sample.wantErr {
				t.Errorf("did not expect an error, got one")
			}
		})
	}
}
