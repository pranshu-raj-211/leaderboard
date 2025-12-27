package unit_test

import (
	"leaderboard/src/backend"
	"testing"
)

func TestAddClient(t *testing.T) {
	store := &fakeLeaderboardStore{}
	cfg := &backend.BroadcasterConfig{
		BroadcastBufferSize:      10,
		PollingIntervalSeconds:   1,
		TopPlayersLimit:          10,
		HeartbeatIntervalSeconds: 5,
	}

	lb, err := backend.CreateLeaderboardBroadcaster(store, cfg)
	if err != nil {
		t.Fatalf("incorrect config passed to leaderboard constructor")
	}
	client, ch := lb.AddClient()

	if client.ID != 1 {
		t.Fatalf("expected client ID to be 1, got %d", client.ID)
	}

	if count := lb.CountClients(); count != 1 {
		t.Fatalf("expected 1 element in clients map, got %d", count)
	}

	select {
	case <-ch:
		t.Fatalf("expected client channel to be empty")
	default:
	}
}

func TestRemoveClient(t *testing.T) {
	store := &fakeLeaderboardStore{}
	cfg := &backend.BroadcasterConfig{
		BroadcastBufferSize:      10,
		PollingIntervalSeconds:   1,
		TopPlayersLimit:          10,
		HeartbeatIntervalSeconds: 5,
	}

	lb, err := backend.CreateLeaderboardBroadcaster(store, cfg)
	if err != nil {
		t.Fatalf("incorrect config passed to leaderboard constructor")
	}
	client, ch := lb.AddClient()

	lb.RemoveClient(client)

	_, ok := <-ch
	if ok {
		t.Fatalf("expected channel to be closed")
	}
	if count := lb.CountClients(); count != 0 {
		t.Fatalf("expected 0 elements in clients map, got %d", count)
	}

	// idempotent, doing it multiple times does not panic
	lb.RemoveClient(client)
}

// create multiple clients, force updates using the BroadcastNow func, check if messages are same

// func TestBroadcast(t *testing.T){

// }
