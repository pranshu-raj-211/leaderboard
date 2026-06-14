package unit_test

import (
	"net/http"
	"testing"
)

func TestRegisterPlayer(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		wantStatus int
		wantStored bool // whether SetPlayerName should have been invoked
	}{
		{
			name:       "valid registration",
			body:       `{"player_id":"p_0001","name":"SwiftFalcon"}`,
			wantStatus: http.StatusOK,
			wantStored: true,
		},
		{
			name:       "missing name",
			body:       `{"player_id":"p_0001"}`,
			wantStatus: http.StatusBadRequest,
			wantStored: false,
		},
		{
			name:       "missing player id",
			body:       `{"name":"SwiftFalcon"}`,
			wantStatus: http.StatusBadRequest,
			wantStored: false,
		},
		{
			name:       "malformed json",
			body:       `{"player_id":`,
			wantStatus: http.StatusBadRequest,
			wantStored: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeLeaderboardStore{}
			r := CreateTestRouter(store)

			w := performRequest(r, "POST", "/players", tc.body)

			if w.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d (%s)", tc.wantStatus, w.Code, w.Body.String())
			}

			gotStored := store.callCount("SetPlayerName") > 0
			if gotStored != tc.wantStored {
				t.Fatalf("expected stored=%v, got %v", tc.wantStored, gotStored)
			}
			if tc.wantStored && store.playerNames["p_0001"] != "SwiftFalcon" {
				t.Fatalf("expected name to be stored, got %q", store.playerNames["p_0001"])
			}
		})
	}
}
