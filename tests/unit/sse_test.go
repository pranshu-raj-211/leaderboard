package unit_test

import (
	"leaderboard/src/backend"
	"testing"
)

func TestBroadcasterCreation(t *testing.T) {
	b := backend.CreateLeaderboardBroadcaster()
	switch v := interface{}(b).(type) {
	case *backend.LeaderboardBroadcaster:
		// correct type
	default:
		t.Errorf("Expected *backend.LeaderboardBroadcaster, got %T", v)
	}
}
