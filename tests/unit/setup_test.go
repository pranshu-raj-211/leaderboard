package unit_test

import (
	"context"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type fakeLeaderboardStore struct {
	timesCalled int
	player1ID   string
	player2ID   string
	result      int

	// configurable fakes
	GetPlayerScoreFn func(ctx context.Context, playerID string) (int64, float64, error)
	GetTopNPlayersFn func(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error)
}

func (store *fakeLeaderboardStore) UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error {
	store.player1ID = player1ID
	store.player2ID = player2ID
	store.timesCalled += 1
	store.result = result
	return nil
}

func (store *fakeLeaderboardStore) GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error) {
	store.timesCalled += 1
	if store.GetPlayerScoreFn != nil {
		return store.GetPlayerScoreFn(ctx, playerID)
	}
	return 0, 0.0, nil
}

func (store *fakeLeaderboardStore) GetTopNPlayers(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
	store.timesCalled += 1
	if store.GetTopNPlayersFn != nil {
		return store.GetTopNPlayersFn(ctx, limit)
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
	r.GET("/player/:id/stats", backend.GetPlayerResults(store))
	r.GET("/leaderboard", backend.GetLeaderboard(store, 10))
	return r
}
