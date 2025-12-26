package unit_test

import (
	"context"
	"encoding/json"
	"errors"
	"leaderboard/src/interfaces"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeaderboard_Success(t *testing.T) {
	store := &fakeLeaderboardStore{
		GetTopNPlayersFn: func(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
			return []interfaces.LeaderboardEntry{
				{PlayerID: "1", Score: 34},
				{PlayerID: "2", Score: 12},
			}, nil
		},
	}
	r := CreateTestRouter(store)

	w := httptest.NewRecorder()
	var resp []interfaces.LeaderboardEntry
	req, _ := http.NewRequest("GET", "/leaderboard", nil)

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

	if len(resp) != 2 {
		t.Fatalf("expected 2 lb entries, got %d", len(resp))
	}

	if resp[0].PlayerID != "1" || resp[0].Score != 34 {
		t.Fatalf("unexpected first entry: %+v", resp[0])
	}

	if resp[1].PlayerID != "2" || resp[1].Score != 12 {
		t.Fatalf("unexpected second entry: %+v", resp[1])
	}
}

func TestLeaderboard_StoreDown(t *testing.T){
	store := &fakeLeaderboardStore{
		GetTopNPlayersFn: func(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
			return nil, errors.New("store down")
		},
	}
	r := CreateTestRouter(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/leaderboard", nil)

	r.ServeHTTP(w, req)

	if w.Code!=http.StatusInternalServerError{
		t.Fatalf("expected code 500, got %d", w.Code)
	}
}