package unit_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestPlayer_Success(t *testing.T) {
	store := &fakeLeaderboardStore{
		playerNames: map[string]string{"12": "SwiftFalcon"},
		GetPlayerScoreFn: func(ctx context.Context, playerID string) (int64, float64, error) {
			return 1, 2.5, nil
		},
	}
	r := CreateTestRouter(store)

	var resp struct {
		PlayerID string  `json:"player_id"`
		Name     string  `json:"name"`
		Rank     int64   `json:"rank"`
		Score    float64 `json:"score"`
	}
	w := performRequest(r, "GET", "/player/12/stats", "")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error while unmarshaling response: %v", err)
	}

	// stats handler resolves the score and the display name exactly once each
	if got := store.callCount("GetPlayerScore"); got != 1 {
		t.Fatalf("expected GetPlayerScore called once, got %d", got)
	}
	if got := store.callCount("GetPlayerName"); got != 1 {
		t.Fatalf("expected GetPlayerName called once, got %d", got)
	}

	if resp.Rank != 1 {
		t.Fatalf("expected rank 1, got %d", resp.Rank)
	}

	if resp.Score != 2.5 {
		t.Fatalf("expected score 2.5, got %f", resp.Score)
	}

	if resp.Name != "SwiftFalcon" {
		t.Fatalf("expected name SwiftFalcon, got %q", resp.Name)
	}
}

func TestPlayer_BadRequest(t *testing.T) {
	store := &fakeLeaderboardStore{
		GetPlayerScoreFn: func(ctx context.Context, playerID string) (int64, float64, error) {
			return 1, 2.5, nil
		},
	}
	r := CreateTestRouter(store)

	w := performRequest(r, "GET", "/player//stats", "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", w.Code)
	}
}

func TestPlayer_StoreDown(t *testing.T) {
	store := &fakeLeaderboardStore{
		GetPlayerScoreFn: func(ctx context.Context, playerID string) (int64, float64, error) {
			return 0, 0.0, errors.New("store is down")
		},
	}
	r := CreateTestRouter(store)

	w := performRequest(r, "GET", "/player/42/stats", "")

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code 500, got %d", w.Code)
	}
}
