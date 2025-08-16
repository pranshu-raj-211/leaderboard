package unit_test

import (
	"leaderboard/src/backend"
	"testing"
)

func TestBroadcasterCreation(t *testing.T) {
	b := backend.CreateLeaderboardBroadcaster()
	defer b.StopBroadcast()
	if _, ok:=any(b).(*backend.LeaderboardBroadcaster); !ok{
		t.Fatalf("Expected *backend.LeaderboardBroadCaster, got %T", b)
	}
}
