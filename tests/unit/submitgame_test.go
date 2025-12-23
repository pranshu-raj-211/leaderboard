package unit_test

import (
	"encoding/json"
	"leaderboard/src/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSubmitGameResults_Success(t *testing.T) {
	store := &fakeSortedSet{}
	r := CreateTestRouter(store)

	w := httptest.NewRecorder()
	gameResult := &models.GameResult{GameID: "12", ServerID: "54", Player1ID: "1", Player2ID: "2", Result: 2}
	gameJSON, _ := json.Marshal(gameResult)
	req, _ := http.NewRequest("POST", "/submit-game", strings.NewReader(string(gameJSON)))
	// req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", w.Code)
	}

	if !store.called {
		t.Fatalf("expected store called to be true")
	}

	if store.player1ID != "1" || store.player2ID != "2" || store.result != 2 {
		t.Fatalf(
			"unexpected store args: p1=%s p2=%s result=%d",
			store.player1ID,
			store.player2ID,
			store.result,
		)
	}
}

type failureCase struct {
	name       string
	gameResult string
}

var incorrectJSON = `{"message":"hello"}`
var missingGameID = models.GameResult{
		GameID:    "",
		Player1ID: "1",
		Player2ID: "2",
		ServerID:  "s1",
		Result:    1,
	}
var missingGameIDJSON, _ = json.Marshal(missingGameID)
var samePlayerID = models.GameResult{
			GameID:    "g1",
			Player1ID: "2",
			Player2ID: "2",
			ServerID:  "s1",
			Result:    1,
		}
var samePlayerIDJSON, _ = json.Marshal(samePlayerID)

var failureTestCases = []failureCase{
	{
		name:       "Non JSON body",
		gameResult: "2",
	},
	{
		name:       "Incorrect JSON - does not match model",
		gameResult: incorrectJSON,
	},
	{
		name: "Invalid GameResult object - GameID missing",
		gameResult: string(missingGameIDJSON),
	},
	{
		name: "Same Player ids",
		gameResult: string(samePlayerIDJSON),
	},
	// too many validation combinations can be tried, not the scope of this test - to be done by test for models instead.
}

func TestSubmitGameResults_Failure(t *testing.T) {
	for _, testCase := range failureTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			store := &fakeSortedSet{}
			r := CreateTestRouter(store)
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("POST", "/submit-game", strings.NewReader(testCase.gameResult))
			r.ServeHTTP(w, req)

			// failure cases - want these requests to fail
			if w.Code == http.StatusOK {
				t.Fatalf("expected status code not equal to 200")
			}

			if store.called {
				t.Fatalf("store should not be called in failure")
			}
		})
	}
}
