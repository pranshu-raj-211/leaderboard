package unit_test

import (
	"context"
	"leaderboard/src/backend"
	"leaderboard/src/config"

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

func init() {
	config.SetLogger(zap.NewNop())
}

func CreateTestRouter(store backend.SortedSetStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.POST("/submit-game", backend.SubmitGameResults(store))
	return r
}
