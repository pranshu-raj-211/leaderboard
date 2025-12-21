package unit_test

import (
	"context"
	"leaderboard/src/backend"
	"leaderboard/src/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type fakeSortedSet struct{}

func (fakeSortedSet) UpdateLeaderboard(ctx context.Context, player1ID string, player2ID string, result int) error {
	return nil
}

func init() {
	config.SetLogger(zap.NewNop())

	r := gin.New()
	store := &fakeSortedSet{}

	r.POST("/submit-game", backend.SubmitGameResults(store))
}
