// Drives synthetic write traffic against the leaderboard.

package simulator

import (
	"context"
	"fmt"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"leaderboard/src/metrics"
	"math/rand"
	"sync"
	"time"
)

type Config struct {
	Enabled            bool
	NumPlayers         int
	TickIntervalMillis int
	MaxConcurrent      int // upper bound on matches played per tick
	Seed               int64
}

type Simulator struct {
	store   interfaces.LeaderboardStore
	cfg     Config
	players []string
	rng     *rand.Rand
	rngMu   sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New(store interfaces.LeaderboardStore, cfg Config) *Simulator {
	ctx, cancel := context.WithCancel(context.Background())
        if cfg.TickIntervalMillis < 0 {
                cfg.TickIntervalMillis = -cfg.TickIntervalMillis
        }
	return &Simulator{
		store:  store,
		cfg:    cfg,
		rng:    rand.New(rand.NewSource(cfg.Seed)),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Seeds the player pool and launches the match loop in the background. Is a no-op if the simulator is disabled.
func (s *Simulator) Start() {
	if !s.cfg.Enabled {
		config.Info("Match simulator disabled", map[string]any{})
		return
	}
	if s.cfg.NumPlayers < 2 {
		config.Error("Simulator needs at least 2 players, not starting", map[string]any{"num_players": s.cfg.NumPlayers})
		return
	}
        if s.cfg.TickIntervalMillis == 0 {
                config.Error("Simulator tick interval must be positive", map[string]any{"tick_ms": s.cfg.TickIntervalMillis})
        }

	s.seedPlayers()

	s.wg.Add(1)
	go s.run()

	config.Info("Match simulator started", map[string]any{
		"num_players":    s.cfg.NumPlayers,
		"tick_ms":        s.cfg.TickIntervalMillis,
		"max_concurrent": s.cfg.MaxConcurrent,
	})
}

// Signals the match loop to exit and waits for in-flight matches to finish (graceful degradation like).
func (s *Simulator) Stop() {
	s.cancel()
	s.wg.Wait()
	config.Info("Match simulator stopped", map[string]any{})
}

// Creates the player pool and registers a display name for each.
func (s *Simulator) seedPlayers() {
	s.players = make([]string, 0, s.cfg.NumPlayers)
	for i := 0; i < s.cfg.NumPlayers; i++ {
		id := fmt.Sprintf("p_%04d", i)
		s.players = append(s.players, id)
		if err := s.store.SetPlayerName(s.ctx, id, randomName(s.rng, i)); err != nil {
			config.Error("Failed to seed player name", map[string]any{"Error": err, "player_id": id})
		}
	}
}

func (s *Simulator) run() {
	defer s.wg.Done()

	interval := time.Duration(s.cfg.TickIntervalMillis) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.playMatches()
		case <-s.ctx.Done():
			return
		}
	}
}

// Plays a random number of matches concurrently.
func (s *Simulator) playMatches() {
	s.rngMu.Lock()
	count := 1 + s.rng.Intn(maxInt(1, s.cfg.MaxConcurrent))
	s.rngMu.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		p1, p2, result := s.randomMatch()
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.store.UpdateLeaderboard(s.ctx, p1, p2, result); err != nil {
				config.Error("Simulated match failed to record", map[string]any{"Error": err})
				return
			}
			metrics.GameSubmissions.Inc()
		}()
	}
	wg.Wait()
}

// Picks two distinct random players and a random result (0/1/2).
func (s *Simulator) randomMatch() (player1ID, player2ID string, result int) {
	s.rngMu.Lock()
	defer s.rngMu.Unlock()

	i := s.rng.Intn(len(s.players))
	j := s.rng.Intn(len(s.players))
	for j == i {
		j = s.rng.Intn(len(s.players))
	}
	return s.players[i], s.players[j], s.rng.Intn(3)
}

var (
	adjectives = []string{"Swift", "Brave", "Silent", "Crimson", "Golden", "Shadow", "Iron", "Mystic", "Frost", "Blazing"}
	nouns      = []string{"Falcon", "Tiger", "Wolf", "Phoenix", "Dragon", "Viper", "Knight", "Ranger", "Specter", "Titan"}
)

func randomName(rng *rand.Rand, n int) string {
	return fmt.Sprintf("%s%s%d", adjectives[rng.Intn(len(adjectives))], nouns[rng.Intn(len(nouns))], n)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
