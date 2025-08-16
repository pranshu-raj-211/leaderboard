package unit_test

import (
	"leaderboard/src/backend"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBroadcasterCreation(t *testing.T) {
	type broadcaster interface {
		StopBroadcast()
		StreamLeaderboard(*gin.Context)
	}
	b := backend.CreateLeaderboardBroadcaster()
	// not adding defer b.StopBroadcast, created nil pointer dereference error, possibly because of Fatalf
	if _, ok := any(b).(broadcaster); !ok {
		t.Fatalf("broadcaster does not satisfy interface requirements, got %T", b)
	}
}
