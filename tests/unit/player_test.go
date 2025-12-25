package unit_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPlayer_Success(t *testing.T) {
	store := &fakeLeaderboardStore{
		GetPlayerScoreFn: func(ctx context.Context, playerID string) (int64, float64, error) {
			return 1, 2.5, nil
		},
	}
	r := CreateTestRouter(store)

	w := httptest.NewRecorder()
	var resp struct {
		Rank  int64
		Score float64
	}
	req, _ := http.NewRequest("GET", "/player/12/stats", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("error while unmarshaling response: %v", err)
	}

	if !(store.timesCalled == 1) {
		t.Fatalf("expected store to be called, timesCalled is %d", store.timesCalled)
	}

	if resp.Rank != 1 {
		t.Fatalf("expected rank 1, got %d", resp.Rank)
	}

	if resp.Score != 2.5 {
		t.Fatalf("expected score 2.5, got %f", resp.Score)
	}
}

func TestPlayer_BadRequest(t *testing.T) {
	store := &fakeLeaderboardStore{
		GetPlayerScoreFn: func(ctx context.Context, playerID string) (int64, float64, error) {
			return 1, 2.5, nil
		},
	}
	r := CreateTestRouter(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/player//stats", nil)

	r.ServeHTTP(w, req)

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

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/player/42/stats", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code 500, got %d", w.Code)
	}
}
