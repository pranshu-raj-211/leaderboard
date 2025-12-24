package interfaces

import (
	"github.com/gin-gonic/gin"
)

type Broadcaster interface {
	StopBroadcast()
	StreamLeaderboard(*gin.Context)
}
