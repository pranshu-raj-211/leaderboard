package simulator

import (
	"context"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func init() {
	config.SetLogger(zap.NewNop())
}

// stubStore is a minimal concurrency-safe LeaderboardStore for simulator tests.
// The simulator plays matches from multiple goroutines, so all access is locked.
type stubStore struct {
	mu      sync.Mutex
	names   map[string]string
	matches int
}

func newStubStore() *stubStore {
	return &stubStore{names: make(map[string]string)}
}

func (s *stubStore) UpdateLeaderboard(ctx context.Context, p1, p2 string, result int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.matches++
	return nil
}

func (s *stubStore) SetPlayerName(ctx context.Context, playerID, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.names[playerID] = name
	return nil
}

func (s *stubStore) nameCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.names)
}

func (s *stubStore) matchCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.matches
}

// Unused by the simulator but required to satisfy the interface.
func (s *stubStore) GetTopNPlayers(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
	return nil, nil
}
func (s *stubStore) GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error) {
	return 0, 0, nil
}
func (s *stubStore) GetPlayerName(ctx context.Context, playerID string) (string, error) {
	return "", nil
}
func (s *stubStore) GetPlayerNames(ctx context.Context, ids []string) (map[string]string, error) {
	return nil, nil
}

func newTestSim(store interfaces.LeaderboardStore, numPlayers, maxConcurrent int) *Simulator {
	return New(store, Config{
		Enabled:            true,
		NumPlayers:         numPlayers,
		TickIntervalMillis: 1,
		MaxConcurrent:      maxConcurrent,
		Seed:               42,
	})
}

func TestSeedPlayers(t *testing.T) {
	store := newStubStore()
	s := newTestSim(store, 10, 3)

	s.seedPlayers()

	if len(s.players) != 10 {
		t.Fatalf("expected 10 players in pool, got %d", len(s.players))
	}
	if got := store.nameCount(); got != 10 {
		t.Fatalf("expected 10 names registered, got %d", got)
	}
}

func TestRandomMatchPicksDistinctPlayers(t *testing.T) {
	store := newStubStore()
	s := newTestSim(store, 5, 3)
	s.seedPlayers()

	for i := 0; i < 1000; i++ {
		p1, p2, result := s.randomMatch()
		if p1 == p2 {
			t.Fatalf("randomMatch returned identical players: %s", p1)
		}
		if result < 0 || result > 2 {
			t.Fatalf("result out of range [0,2]: %d", result)
		}
	}
}

func TestPlayMatchesRespectsBound(t *testing.T) {
	store := newStubStore()
	const maxConcurrent = 5
	s := newTestSim(store, 8, maxConcurrent)
	s.seedPlayers()

	// A single batch must play between 1 and maxConcurrent matches.
	s.playMatches()

	got := store.matchCount()
	if got < 1 || got > maxConcurrent {
		t.Fatalf("expected between 1 and %d matches, got %d", maxConcurrent, got)
	}
}

func TestStartDisabledIsNoop(t *testing.T) {
	store := newStubStore()
	s := New(store, Config{Enabled: false, NumPlayers: 10, MaxConcurrent: 3, TickIntervalMillis: 1})

	s.Start()
	defer s.Stop()

	if len(s.players) != 0 {
		t.Fatalf("disabled simulator should not seed players, got %d", len(s.players))
	}
	if got := store.matchCount(); got != 0 {
		t.Fatalf("disabled simulator should not play matches, got %d", got)
	}
}

func TestStartStopProducesWrites(t *testing.T) {
	store := newStubStore()
	s := newTestSim(store, 6, 3)

	s.Start()
	// Stop blocks until the loop exits; with a 1ms tick at least one batch runs.
	s.Stop()

	if got := store.nameCount(); got != 6 {
		t.Fatalf("expected 6 players seeded, got %d", got)
	}
	if store.matchCount() == 0 {
		t.Fatalf("expected at least one match to be recorded after start/stop")
	}
}
