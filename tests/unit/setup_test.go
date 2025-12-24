package unit_test

import (
	"context"
	"leaderboard/src/backend"
	"leaderboard/src/config"
	"leaderboard/src/interfaces"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type fakeSortedSet struct {
	called    bool
	player1ID string
	player2ID string
	result    int
}

func (store *fakeSortedSet) UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error {
	store.player1ID = player1ID
	store.player2ID = player2ID
	store.called = true
	store.result = result
	return nil
}

func (store *fakeSortedSet) GetPlayerScore(ctx context.Context, playerID string) (int64, float64, error) {
	return 1, 1.0, nil
}

func (store *fakeSortedSet) GetTopNPlayers(ctx context.Context, limit int64) ([]interfaces.LeaderboardEntry, error) {
	// TODO: implement good mocks
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
	r.GET("/leaderboard", backend.GetLeaderboard(store))
	return r
}
