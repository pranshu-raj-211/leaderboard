package unit_test

import (
	"go.uber.org/zap"
	"leaderboard/src/config"
)

func init() {
	config.SetLogger(zap.NewNop())
}
