package unit_test

import (
	"context"
	"io"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// fakeLeaderboardStore is a concurrency-safe in-memory test double. It records both a total call count (timesCalled, kept for existing assertions) and a per-method count so tests can assert on a specific method without coupling to the total number of store interactions.
type fakeLeaderboardStore struct {
	mu        sync.Mutex
	calls     map[string]int
	player1ID string
	player2ID string
	result    int

	playerNames map[string]string

	// configurable fakes
	GetPlayerScoreFn func(ctx context.Context, playerID string) (int64, float64, error)
	GetTopNPlayersFn func(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error)
}

// Increments the per-method and total counters under the lock.
func (store *fakeLeaderboardStore) record(method string) {
	if store.calls == nil {
		store.calls = make(map[string]int)
	}
	store.calls[method]++
}

// Returns how many times a specific method was invoked.
func (store *fakeLeaderboardStore) callCount(method string) int {
	store.mu.Lock()
	defer store.mu.Unlock()
	return store.calls[method]
}

// Returns the number of store method invocations across all methods.
func (store *fakeLeaderboardStore) totalCalls() int {
	store.mu.Lock()
	defer store.mu.Unlock()
	total := 0
	for _, n := range store.calls {
		total += n
	}
	return total
}

func (store *fakeLeaderboardStore) SetPlayerName(ctx context.Context, playerID string, name string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.record("SetPlayerName")
	if store.playerNames == nil {
		store.playerNames = make(map[string]string)
	}
	store.playerNames[playerID] = name
	return nil
}

func (store *fakeLeaderboardStore) GetPlayerName(ctx context.Context, playerID string) (string, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.record("GetPlayerName")
	return store.playerNames[playerID], nil
}

func (store *fakeLeaderboardStore) GetPlayerNames(ctx context.Context, playerIDs []string) (map[string]string, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.record("GetPlayerNames")
	names := make(map[string]string, len(playerIDs))
	for _, id := range playerIDs {
		if name, ok := store.playerNames[id]; ok {
			names[id] = name
		}
	}
	return names, nil
}

func (store *fakeLeaderboardStore) UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.record("UpdateLeaderboard")
	store.player1ID = player1ID
	store.player2ID = player2ID
	store.result = result
	return nil
}

func (store *fakeLeaderboardStore) GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error) {
	store.mu.Lock()
	store.record("GetPlayerScore")
	fn := store.GetPlayerScoreFn
	store.mu.Unlock()
	if fn != nil {
		return fn(ctx, playerID)
	}
	return 0, 0.0, nil
}

func (store *fakeLeaderboardStore) GetTopNPlayers(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
	store.mu.Lock()
	store.record("GetTopNPlayers")
	fn := store.GetTopNPlayersFn
	store.mu.Unlock()
	if fn != nil {
		return fn(ctx, limit)
	}
	return nil, nil
}

func init() {
	config.SetLogger(zap.NewNop())
}

func CreateTestRouter(store interfaces.LeaderboardStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/submit-game", backend.SubmitGameResults(store))
	r.POST("/players", backend.RegisterPlayer(store))
	r.GET("/player/:id/stats", backend.GetPlayerResults(store))
	r.GET("/leaderboard", backend.GetLeaderboard(store, 10))
	return r
}

// Builds a request (with an optional body) and serves it against the router, returning the recorder. body may be "" for GET requests.
func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, _ := http.NewRequest(method, path, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
